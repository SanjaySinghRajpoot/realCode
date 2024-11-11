from django.shortcuts import render
from django.http import HttpResponse
from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
from django.core.cache import cache
import subprocess
import json
import os
import requests
from dotenv import load_dotenv

load_dotenv()

api_key = os.getenv('google_api_key')

# Create your views here.
def homepage(request):
    return HttpResponse("Welcome to Real code!!")

def generate_content(code):
    try: 
        url = f"https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key={api_key}"
        headers = {
            'Content-Type': 'application/json',
        }
        data = {
            "contents": [{
                "parts": [{"text": f"Evaluate the following code and provide feedback on optimization, code quality, code structure, and reliability. Give marks out of 10 for each category. Summarize it in 100 words\n\n{code}"}]
            }]
        }
        response = requests.post(url, headers=headers, json=data)
        return response.json()
    except Exception as e:
        return {'error': str(e)}

@csrf_exempt 
def compile_code(request):

    if request.method == 'POST':
        data = json.loads(request.body)
        code_txt = data['code']
        language = data['language']


        result = cache.get(code_txt)
        if result:
            return JsonResponse({'output': result})

        if language == 'python':
            try:
                # Execute the Python code
                result = subprocess.run(
                    ['python3', '-c', code_txt],
                    capture_output=True,
                    text=True,
                    check=True
                )

                feedback = generate_content(code_txt)

                res = feedback.get('candidates')

                cache.set(code_txt, {'output': result.stdout, 'feedback': res[0]}, timeout=60)

                # Return the output to the user
                return JsonResponse({'output': result.stdout, 'feedback': res[0]})
            except subprocess.CalledProcessError as e:
                # Return the error to the user
                return JsonResponse({'error': e.stderr}, status=400)
        
        # Handle other languages if needed
        return JsonResponse({'error': 'Unsupported language'}, status=400)

    return JsonResponse({'error': 'Invalid request method'}, status=405)





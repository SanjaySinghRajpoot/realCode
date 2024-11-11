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
    
def compile_and_run_golang(code_txt):
    """Compiles and runs Golang code, returning output and feedback.

    Args:
        code_txt (str): The Golang code to compile and execute.

    Returns:
        dict: A dictionary containing the following keys:
            output (str, optional): The standard output of the compiled and executed code.
            error (str, optional): The standard error stream of the compilation or execution process.
            feedback (dict, optional): Feedback generated from the code (implementation-specific).
    """

    # Create a temporary file to store the Golang code
    with open("temp_golang.go", "w") as f:
        f.write(code_txt)

    # Compile the Golang code
    try:
        compilation_process = subprocess.run(
            ["go", "build", "-o", "temp_golang", "temp_golang.go"],
            check=True,
            capture_output=True,
            text=True,  # Ensure text output for easier handling
        )
    except subprocess.CalledProcessError as e:
        return {"error": e.stderr}  # Return compilation error

    # Execute the compiled program
    try:
        execution_process = subprocess.run(
            ["./temp_golang"],
            capture_output=True,
            text=True,  # Ensure text output
            timeout=5,  # Set a timeout to prevent infinite execution (adjust as needed)
        )
    except subprocess.CalledProcessError as e:
        return {"error": e.stderr}  # Return execution error
    except subprocess.TimeoutExpired as e:
        return {"error": "Golang program execution timed out"}

    # Remove the temporary file
    try:
        subprocess.run(["rm", "temp_golang.go", "temp_golang"], check=True)
    except subprocess.CalledProcessError:
        pass  # Ignore errors during cleanup

    # Process the output and generate feedback (implementation-specific)
    output = execution_process.stdout
    feedback = generate_content(code_txt)  # Replace with your feedback generation logic

    return {"output": output, "feedback": feedback.get('candidates')[0] if feedback else None}

@csrf_exempt 
def compile_code(request):

    if request.method == 'POST':
        data = json.loads(request.body)
        code_txt = data['code']
        language = data['language']


        result = cache.get(code_txt)
        if result:
            return JsonResponse(result)

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
            
        elif language == 'golang':
            result = compile_and_run_golang(code_txt)

            cache.set(code_txt, result, timeout=60)

            return JsonResponse(result, status=400 if 'error' in result else 200)
        
        # Handle other languages if needed
        return JsonResponse({'error': 'Unsupported language'}, status=400)

    return JsonResponse({'error': 'Invalid request method'}, status=405)





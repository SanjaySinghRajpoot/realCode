import { useState } from "react";
import { Box, Button, Text, useToast } from "@chakra-ui/react";
import { executeCode } from "../api";
import ReactMarkdown from "react-markdown";

const Output = ({ editorRef, language }) => {
  const toast = useToast();
  const [output, setOutput] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [isError, setIsError] = useState(false);
  const [feedback, setFeedback] = useState(null);

  const runCode = async () => {
    const sourceCode = editorRef.current.getValue();
    if (!sourceCode) return;
    try {
      setIsLoading(true);
      const { output, feedback } = await executeCode(language, sourceCode);
      setOutput(output.split("\n"));
      const feedbackText = feedback.content.parts[0].text;
      setFeedback(feedbackText);
      // check for output error code
      setIsError(false)
      // data.stderr ? setIsError(true) : setIsError(false);
    } catch (error) {
      console.log(error);
      toast({
        title: "An error occurred.",
        description: error.message || "Unable to run code",
        status: "error",
        duration: 6000,
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Box w="50%">
      {/* <Text mb={2} fontSize="lg">
        Output
      </Text> */}
      <Button
        variant="outline"
        colorScheme="green"
        mb={4}
        isLoading={isLoading}
        onClick={runCode}
      >
        Run Code
      </Button>
      <Box
        height="38vh"
        p={2}
        color={isError ? "red.400" : ""}
        border="1px solid"
        borderRadius={4}
        borderColor={isError ? "red.500" : "#333"}
        overflowY="auto" 
      >
        {output
          ? output.map((line, i) => <Text key={i}>{line}</Text>)
          : 'Click "Run Code" to see the output here'}
      </Box>

      <Text mb={2} mt={4} fontSize="lg">
        Code Analysis
      </Text>
      <Box
        height="28vh"
        p={2}
        color={isError ? "red.400" : ""}
        border="1px solid"
        borderRadius={4}
        borderColor={isError ? "red.500" : "#333"}
        overflowY="auto" 
      >
        {feedback
          ? <ReactMarkdown>{feedback}</ReactMarkdown>
          : 'Code analysis will be displayed here'}
      </Box>
    </Box>
  );
};
export default Output;

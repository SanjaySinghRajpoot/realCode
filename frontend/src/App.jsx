import { Box, Text } from "@chakra-ui/react";
import CodeEditor from "./components/CodeEditor";

function App() {
  return (
    <Box minH="100vh" bg="#0f0a19" color="gray.500" px={6} py={8} overflow="hidden">
      <Text mb={2} ml={1} fontSize="2xl" fontWeight="bold" textAlign="left" color="white">
        RealCode
      </Text>
      <CodeEditor />
    </Box>
  );
}

export default App;

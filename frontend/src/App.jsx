import { Box, Text } from "@chakra-ui/react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import CodeEditor from "./components/CodeEditor";
import GoogleLoginPage from "./components/GoogleLoginPage";

function App() {
  return (
    <Router>
      <Box minH="100vh" bg="#0f0a19" color="gray.500" px={[4, 6, 8]} py={[4, 6, 8]} overflow="hidden">
        <Text mb={[1, 2]} ml={[0, 1]} fontSize={["xl", "2xl"]} fontWeight="bold" textAlign="left" color="white">
          RealCode
        </Text>
        <Routes>
          <Route path="/login" element={<GoogleLoginPage />} />
          <Route path="/" element={<CodeEditor />} />
        </Routes>
      </Box>
    </Router>
  );
}

export default App;

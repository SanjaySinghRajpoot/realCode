import React from "react";
import { GoogleOAuthProvider, GoogleLogin } from '@react-oauth/google';
import { Box, Text } from "@chakra-ui/react";
import { useNavigate } from "react-router-dom";

const GoogleLoginPage = () => {
  const navigate = useNavigate();

  const responseGoogle = (response) => {
    console.log(response);
    // Assuming successful login, redirect to the code editor path
    navigate("/");
  };

  return (
    <GoogleOAuthProvider clientId="287564293735-sdfqtpf6vlt6evkvqb5o4ink090voekf.apps.googleusercontent.com">
      <Box minH="70vh" bg="#0f0a19" color="gray.500" display="flex" alignItems="center" justifyContent="center">
        <Box textAlign="center">
          <Text mb={4} fontSize="2xl" fontWeight="bold" color="white">
            Login with Google
          </Text>
          <GoogleLogin
            onSuccess={responseGoogle}
            onError={responseGoogle}
          />
        </Box>
      </Box>
    </GoogleOAuthProvider>
  );
};

export default GoogleLoginPage;
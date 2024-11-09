import axios from "axios";
import { LANGUAGE_VERSIONS } from "./constants";

console.log(process.env.REACT_APP_BACKEND_URL)

const API = axios.create({
  baseURL: process.env.BACKEND_URL,
  headers: {
    'Access-Control-Allow-Origin': '*',
    'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
    'Access-Control-Allow-Headers': 'Origin, Content-Type, Accept, Authorization',
  },
});

export const executeCode = async (language, sourceCode) => {
  const response = await API.post("/compile", {
    language: language,
    code: sourceCode
    // version: LANGUAGE_VERSIONS[language],
    // files: [
    //   {
    //     content: sourceCode,
    //   },
    // ],
  });


  return response.data;
};

import axios from "axios";
import { LANGUAGE_VERSIONS } from "./constants";

const API = axios.create({
  baseURL: "http://localhost:8080",
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

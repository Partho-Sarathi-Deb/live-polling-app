import Login from "./pages/Login";
import Signup from "./pages/Signup";
import CreatePoll from "./pages/CreatePoll";
import VotePoll from "./pages/VotePoll";
import "./App.css";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Navigate to="/login" replace />} />
        <Route path="/login" element={<Login />} />
        <Route path="/signup" element={<Signup />} />
        <Route path="/create" element={<CreatePoll />} />
        <Route path="/poll/:id" element={<VotePoll />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Login from "./pages/Login";
import Signup from "./pages/Signup";
import CreatePoll from "./pages/CreatePoll";
import VotePoll from "./pages/VotePoll";
import "./App.css";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/signup" element={<Signup />} />
        <Route path="/create" element={<CreatePoll />} />
        <Route path="/poll/:id" element={<VotePoll />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
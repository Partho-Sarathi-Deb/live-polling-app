import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";

function CreatePoll() {
  const [question, setQuestion] = useState("");
  const [options, setOptions] = useState(["", ""]);
  const [error, setError] = useState("");
  const navigate = useNavigate();
    useEffect(() => {
    if (!localStorage.getItem("token")) {
      navigate("/login");
    }
  }, [navigate]);

  const updateOption = (index, value) => {
    const newOptions = [...options];
    newOptions[index] = value;
    setOptions(newOptions);
  };

  const addOption = () => setOptions([...options, ""]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");

    const token = localStorage.getItem("token");
    if (!token) {
      navigate("/login");
      return;
    }

    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/polls`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ question, options: options.filter((o) => o.trim() !== "") }),
      });

      const data = await res.json();

      if (!res.ok) {
        setError(data.error || "Could not create poll");
        return;
      }

      navigate(`/poll/${data.id}`);
    } catch (err) {
      setError("Could not reach server");
    }
  };

  return (
    <div>
      <h2>Create Poll</h2>
      <form onSubmit={handleSubmit}>
        <input
          placeholder="Question"
          value={question}
          onChange={(e) => setQuestion(e.target.value)}
        />
        {options.map((opt, i) => (
          <input
            key={i}
            placeholder={`Option ${i + 1}`}
            value={opt}
            onChange={(e) => updateOption(i, e.target.value)}
          />
        ))}
        <button type="button" onClick={addOption}>Add option</button>
        <button type="submit">Create Poll</button>
      </form>
      {error && <p className="error">{error}</p>}
    </div>
  );
}

export default CreatePoll;
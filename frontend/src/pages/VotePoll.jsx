import { useState, useEffect } from "react";
import { useParams } from "react-router-dom";

function VotePoll() {
  const { id } = useParams();
  const [poll, setPoll] = useState(null);
  const [counts, setCounts] = useState({});
  const [voted, setVoted] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    fetch(`http://localhost:8080/polls/${id}`)
      .then((res) => res.json())
      .then((data) => setPoll(data))
      .catch(() => setError("Could not load poll"));
  }, [id]);

  useEffect(() => {
    const ws = new WebSocket(`${import.meta.env.VITE_WS_URL}/ws/polls/${id}`);
    ws.onmessage = (event) => {
      setCounts(JSON.parse(event.data));
    };
    return () => ws.close();
  }, [id]);

  const handleVote = async (optionIndex) => {
    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/polls/${id}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ optionIndex }),
      });
      if (!res.ok) {
        const data = await res.json();
        setError(data.error || "Vote failed");
        return;
      }
      setVoted(true);
    } catch (err) {
      setError("Could not reach server");
    }
  };

  if (!poll) return <div>Loading...</div>;

  return (
    <div className="page">
      <h2>{poll.question}</h2>
      {!voted &&
        poll.options.map((opt, i) => (
          <button key={i} onClick={() => handleVote(i)}>
            {opt}
          </button>
        ))}
      {error && <p className="error">{error}</p>}
      <h3>Live Results</h3>
      <ul className="results-list">
        {poll.options.map((opt, i) => (
          <li key={i}>
            {opt}: {counts[i] || 0}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default VotePoll;
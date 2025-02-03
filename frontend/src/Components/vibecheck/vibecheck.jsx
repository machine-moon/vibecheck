import { useState } from "react";
import "./vibecheck.css";

const API_HOST = import.meta.env.API_HOST || "localhost";
const API_PORT = import.meta.env.API_PORT || "8080";
const API_BASE_URL = `http://${API_HOST}:${API_PORT}`;

const VibeCheck = () => {
  const [problem, setProblem] = useState(null);
  const [problemId, setProblemId] = useState("");
  const [ProblemText, setProblemText] = useState("");

  const [hint, setHint] = useState("");
  const [selectedAnswer, setSelectedAnswer] = useState("");
  const [correctAnswers, setCorrectAnswers] = useState(0);
  const [totalAnswers, setTotalAnswers] = useState(0);
  const [submitted, setSubmitted] = useState(false);
  const [isCorrect, setIsCorrect] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchProblem = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await fetch(`${API_BASE_URL}/problem/quiz`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });
      if (!response.ok) throw new Error("Failed to fetch problem");
      const data = await response.json();
      setProblem(data);
      console.log(data);
      setProblemId(data.problem.id);
      console.log(data.problem.id);
      setProblemText(data.problem.text);
      console.log(data.problem.text);
      setHint("");
      setSelectedAnswer("");
      setSubmitted(false);
      setIsCorrect(false);
    } catch (error) {
      console.error("Error fetching problem:", error);
      setError("Failed to fetch problem. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  const fetchHint = async () => {
    if (!problem) return;
    try {
      setLoading(true);
      setError(null);
      const response = await fetch(
        `${API_BASE_URL}/problem/hint/${problemId}`,
        {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
          },
        }
      );
      if (!response.ok) throw new Error("Failed to fetch hint");
      const data = await response.json();
      if (
        !data.hint ||
        data.hint.trim() === "" ||
        data.hint === problem.text ||
        data.hint.toString().toLowerCase() === ProblemText.toLowerCase()
      ) {
        setHint(
          "Hmm, it seems the hint is playing hide and seek. Try solving the problem without it!"
        );
      } else {
        setHint(data.hint);
      }
    } catch (error) {
      console.error("Error fetching hint:", error);
      setError("Failed to fetch hint. Please try again.");
    } finally {
      setLoading(false);
    }
  };
  const submitAnswer = async () => {
    if (!problem || !selectedAnswer) return;
    try {
      setLoading(true);
      setError(null);
      const response = await fetch(`${API_BASE_URL}/problem/answer`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ id: problemId, guess: selectedAnswer }),
      });
      if (!response.ok) throw new Error("Failed to submit answer");
      const data = await response.json();
      setSubmitted(true);
      setIsCorrect(data.correct);
      setTotalAnswers(totalAnswers + 1);
      if (data.correct) {
        setCorrectAnswers(correctAnswers + 1);
      }
    } catch (error) {
      console.error("Error submitting answer:", error);
      setError("Failed to submit answer. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="vibecheck">
      {error && <div className="error">{error}</div>}
      <button onClick={fetchProblem} disabled={loading || submitted}>
        {loading ? "Loading..." : "Skip"}
      </button>
      <button onClick={fetchHint} disabled={!problemId || loading}>
        {loading ? "Loading..." : "Get Hint"}
      </button>
      <div className="problem">{ProblemText || ""}</div>
      <div className="hint majestic-hint">{hint}</div>
      <div className="answers">
        <button
          className={`answer-btn ${
            selectedAnswer === "negative" ? "selected" : ""
          } ${
            submitted && selectedAnswer === "negative"
              ? isCorrect
                ? "correct"
                : "incorrect"
              : ""
          }`}
          onClick={() => setSelectedAnswer("negative")}
          disabled={submitted || loading}
        >
          Negative
        </button>
        <button
          className={`answer-btn ${
            selectedAnswer === "neutral" ? "selected" : ""
          } ${
            submitted && selectedAnswer === "neutral"
              ? isCorrect
                ? "correct"
                : "incorrect"
              : ""
          }`}
          onClick={() => setSelectedAnswer("neutral")}
          disabled={submitted || loading}
        >
          Neutral
        </button>
        <button
          className={`answer-btn ${
            selectedAnswer === "positive" ? "selected" : ""
          } ${
            submitted && selectedAnswer === "positive"
              ? isCorrect
                ? "correct"
                : "incorrect"
              : ""
          }`}
          onClick={() => setSelectedAnswer("positive")}
          disabled={submitted || loading}
        >
          Positive
        </button>
      </div>
      {!submitted && (
        <button onClick={submitAnswer} disabled={!selectedAnswer || loading}>
          {loading ? "Loading..." : "Submit"}
        </button>
      )}
      {submitted && <button onClick={fetchProblem}>Next Question</button>}
      <div className="counter">
        Correct: {correctAnswers} / {totalAnswers}
      </div>
    </div>
  );
};

export default VibeCheck;

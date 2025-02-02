"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { getRandomTweet, getHint, checkAnswer } from "./services/api";
import { Alert } from "@/components/ui/alert";

interface Tweet {
  id: string;
  text: string;
  hint?: string;
}

export default function Home() {
  const [tweet, setTweet] = useState<Tweet | null>(null);
  const [hint, setHint] = useState("");
  const [score, setScore] = useState(0);
  const [selectedAnswer, setSelectedAnswer] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const getNewProblem = async () => {
    try {
      setLoading(true);
      setError(null);
      const newTweet = await getRandomTweet();
      setTweet(newTweet);
      setHint("");
      setSelectedAnswer(null);
    } catch (err) {
      console.error("Error fetching new problem:", err);
      setError("Failed to fetch new problem. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  const showHint = async () => {
    if (!tweet) return;

    try {
      setLoading(true);
      setError(null);
      const hintText = await getHint(tweet.id);
      setHint(hintText);
    } catch (err) {
      console.error("Error fetching hint:", err);
      setError("Failed to fetch hint. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  const handleAnswer = (answer: string) => {
    setSelectedAnswer(answer);
  };

  const handleSubmit = async () => {
    if (!tweet || !selectedAnswer) return;

    try {
      setLoading(true);
      setError(null);
      const isCorrect = await checkAnswer(tweet.id, selectedAnswer);
      if (isCorrect) {
        setScore((prev) => prev + 1);
      }
      await getNewProblem();
    } catch (err) {
      console.error("Error submitting answer:", err);
      setError("Failed to submit answer. Please try again.");
      setLoading(false);
    }
  };

  return (
    <main className="min-h-screen bg-gradient-to-b from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800 p-8">
      <div className="max-w-2xl mx-auto space-y-8">
        <div className="text-center space-y-2">
          <h1 className="text-4xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-500 to-teal-500">
            Vibe Check
          </h1>
          <div className="inline-block px-4 py-2 bg-white dark:bg-gray-800 rounded-full shadow-sm">
            <p className="text-xl font-semibold">
              Score: <span className="text-blue-500">{score}</span>
            </p>
          </div>
        </div>

        {error && (
          <Alert variant="destructive" className="animate-fadeIn">
            {error}
          </Alert>
        )}

        <Card className="border-2 shadow-lg">
          <CardContent className="p-6 space-y-6">
            <Button
              onClick={getNewProblem}
              className="w-full bg-gradient-to-r from-blue-500 to-teal-500 hover:from-blue-600 hover:to-teal-600 text-white shadow-md"
              disabled={loading}
            >
              {tweet ? "New Problem" : "Start"}
            </Button>

            {tweet && (
              <div className="space-y-6 animate-fadeIn">
                <div className="p-4 bg-white dark:bg-gray-800 rounded-lg shadow-inner">
                  <p className="text-lg text-center">{tweet.text}</p>
                </div>

                <Button
                  onClick={showHint}
                  variant="outline"
                  className="w-full border-2 hover:bg-gray-50 dark:hover:bg-gray-800"
                  disabled={loading}
                >
                  Show Hint
                </Button>

                {hint && (
                  <div className="p-3 bg-blue-50 dark:bg-gray-800 border-2 border-blue-200 dark:border-gray-700 rounded-lg animate-fadeIn">
                    <p className="text-sm text-center text-blue-600 dark:text-blue-400">
                      💡 {hint}
                    </p>
                  </div>
                )}

                <div className="grid grid-cols-3 gap-3">
                  <Button
                    onClick={() => handleAnswer("positive")}
                    variant={
                      selectedAnswer === "positive" ? "default" : "outline"
                    }
                    className={`${
                      selectedAnswer === "positive"
                        ? "bg-green-500 hover:bg-green-600"
                        : "hover:bg-green-50 dark:hover:bg-gray-800"
                    } border-2`}
                    disabled={loading}
                  >
                    Positive 😊
                  </Button>
                  <Button
                    onClick={() => handleAnswer("neutral")}
                    variant={
                      selectedAnswer === "neutral" ? "default" : "outline"
                    }
                    className={`${
                      selectedAnswer === "neutral"
                        ? "bg-yellow-500 hover:bg-yellow-600"
                        : "hover:bg-yellow-50 dark:hover:bg-gray-800"
                    } border-2`}
                    disabled={loading}
                  >
                    Neutral 😐
                  </Button>
                  <Button
                    onClick={() => handleAnswer("negative")}
                    variant={
                      selectedAnswer === "negative" ? "default" : "outline"
                    }
                    className={`${
                      selectedAnswer === "negative"
                        ? "bg-red-500 hover:bg-red-600"
                        : "hover:bg-red-50 dark:hover:bg-gray-800"
                    } border-2`}
                    disabled={loading}
                  >
                    Negative 😔
                  </Button>
                </div>

                <Button
                  onClick={handleSubmit}
                  disabled={!selectedAnswer || loading}
                  className={`w-full shadow-md ${
                    selectedAnswer && !loading
                      ? "bg-gradient-to-r from-blue-500 to-teal-500 hover:from-blue-600 hover:to-teal-600"
                      : "opacity-50 cursor-not-allowed"
                  }`}
                >
                  {loading ? "Loading..." : "Submit"}
                </Button>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </main>
  );
}

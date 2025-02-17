interface Tweet {
  id: string
  text: string
  hint?: string
  sentiment?: 'positive' | 'negative' | 'neutral'
}

interface Data {
  problem: Problem
}

interface Problem {
  id: string
  text: string
}

interface AttemptSolution {
  id: string
  guess: string
}



const API_HOST = process.env.API_HOST || "localhost";
const API_PORT = process.env.API_PORT || "8080";
const API_BASE_URL = `http://${API_HOST}:${API_PORT}`;



export async function getRandomTweet(): Promise<Tweet> {
  try {
    const response = await fetch(`${API_BASE_URL}/problem/quiz`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });
    if (!response.ok) throw new Error('API error')
    const data: Data = await response.json()
    return {
      id: data.problem.id,
      text: data.problem.text
    }
  } catch (error) {
    console.error('Failed to fetch problem:', error)
    throw error
  }
}

export async function getHint(tweetId: string): Promise<string> {
  try {
    const response = await fetch(
      `${API_BASE_URL}/problem/hint/${tweetId}`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      }
    );
    if (!response.ok) throw new Error('API error')
    const data = await response.json()
    return data.hint
  } catch (error) {
    console.error('Failed to fetch hint:', error)
    throw error
  }
}

export async function checkAnswer(tweetId: string, guess: string): Promise<boolean> {
  try {
    const response = await fetch(`${API_BASE_URL}/problem/answer`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        id: tweetId,
        guess: guess
      } as AttemptSolution)
    })
    if (!response.ok) throw new Error('API error')
    const data = await response.json()
    return data.correct
  } catch (error) {
    console.error('Failed to check answer:', error)
    throw error
  }
} 
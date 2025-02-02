interface Tweet {
  id: string
  text: string
  hint?: string
  sentiment?: 'positive' | 'negative' | 'neutral'
}

interface Problem {
  id: string
  text: string
}

interface AttemptSolution {
  id: string
  guess: string
}

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

// Mock data for when API is not available
const MOCK_TWEETS: Tweet[] = [
  {
    text: "Just had the most amazing coffee at this new cafe! ☕️✨",
    hint: "Think about the emotional words and emojis used",
    sentiment: 'positive'
  },
  {
    text: "The weather is so gloomy today, I don't want to get out of bed 😔",
    hint: "Consider the mood expressed and the emoji used",
    sentiment: 'negative'
  },
  {
    text: "Just finished reading the news for today.",
    hint: "Look for presence or absence of emotional words",
    sentiment: 'neutral'
  }
]

export async function getRandomTweet(): Promise<Tweet> {
  try {
    const response = await fetch(`${API_BASE_URL}/problem/quiz`)
    if (!response.ok) throw new Error('API error')
    const problem: Problem = await response.json()
    return {
      id: problem.id,
      text: problem.text
    }
  } catch (error) {
    console.error('Failed to fetch problem:', error)
    throw error
  }
}

export async function getHint(tweetId: string): Promise<string> {
  try {
    const response = await fetch(`${API_BASE_URL}/problem/hint/${tweetId}`)
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
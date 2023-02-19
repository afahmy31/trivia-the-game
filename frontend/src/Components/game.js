import { useState } from 'react';
import Button from '@mui/material/Button';
import { questions } from './questions.js'
import { Container } from '@mui/system';
// import image from './image.jpg'


export default function Game() {
  const [teamOneScore, setTeamOneScore] = useState(0);
  const [teamTwoScore, setTeamTwoScore] = useState(0);
  const [currentQuestionId, setCurrentQuestionId] = useState(localStorage.getItem("currentQuestionId"));
  const renderScores = () => {
    return (
      <div className='scores'>
        <div className='single-score'>
          <h2 className='score-team-name'> {localStorage.getItem("teamOneName") + " : " + teamOneScore}</h2>
          <div className='score-buttons'>
            <Button className='score-plus-btn' variant='contained' onClick={() => setTeamOneScore(teamOneScore + 1)}>+</Button>
            <Button variant='contained' color='secondary' onClick={() => setTeamOneScore(teamOneScore - 1)}>-</Button>

          </div>
        </div>
        <div className='single-score'>
          <h2 className='score-team-name'> {localStorage.getItem("teamTwoName") + " : " + teamTwoScore}</h2>
          <div className='score-buttons'>
            <Button className='score-plus-btn' variant='contained' onClick={() => setTeamTwoScore(teamTwoScore + 1)}>+</Button>
            <Button variant='contained' color='secondary' onClick={() => setTeamTwoScore(teamTwoScore - 1)}>-</Button>
          </div>
        </div>

      </div>
    )
  }
  const handleNextQuestion = () => {
    setCurrentQuestionId(parseInt(currentQuestionId) + 1)
  }
  const renderQuestion = () => {
    let question = questions[currentQuestionId - 1]
    return (
      <div>
        {renderScores()}
        <h2> {"Question " + question.id}</h2>
        <hr />
        <h1> {question.question}</h1>
        <div className='answers'>
          {question.answers.map((q, i) => {
            return (
              <div className='answer' key={i}>{q.value}</div>
            )
          })}
        </div>
      </div>
    )
  }

  return (
    <Container>
      <div className='game'>
        {renderQuestion()}
        <Button className='next-question' onClick={() => handleNextQuestion()} >Next Question</Button>
        {/* <img src={image} alt="Logo" /> */}
      </div>
    </Container>
  )
}
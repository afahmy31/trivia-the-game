import { useState } from 'react';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import Container from '@mui/material/Container';
import { useNavigate } from "react-router-dom";
import image from './image.jpg'


export default function MainMenu() {
  const navigate = useNavigate();
  const [teamOneName, setTeamOneName] = useState("");
  const [teamTwoName, setTeamTwoName] = useState("");

  const handleNewGame = () => {
    localStorage.setItem("teamOneName", teamOneName)
    localStorage.setItem("teamTwoName", teamTwoName)
    localStorage.setItem("teamOneScore", 0)
    localStorage.setItem("teamTwoScore", 0)
    localStorage.setItem("currentQuestionId", 1)
    navigate("/game");
  }

  return (
    <div className='main-menu'>
      <img className='image' src={image} alt="Logo" />
      <h2>Dodo's Trivia Night</h2>
      <TextField className='name-input' label={"Team One Name"} onChange={(e) => setTeamOneName(e.target.value)}></TextField>
      <TextField className='name-input' label={"Team Two Name"} onChange={(e) => setTeamTwoName(e.target.value)}></TextField>
      <Button variant='contained' onClick={() => handleNewGame()}>New Game</Button>
    </div>
  )
}
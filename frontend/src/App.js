import './App.css';
import {
  BrowserRouter as Router,
  Routes,
  Route,
} from "react-router-dom";
import MainMenu from './Components/mainMenu';
import Game from './Components/game';
import { createTheme, ThemeProvider } from '@mui/material';

const theme = createTheme({
  palette: {
    primary: {
      main: "#843D80"
    },
    secondary: {
      main: "#E3E5E8"
    }
  }
});

export default function App() {
  return (
    <ThemeProvider theme={theme}>
      <Router>
        <Routes>
          <Route path="/" element={<MainMenu />} />
          <Route path="/game" element={<Game />} />
        </Routes>
      </Router>
    </ThemeProvider>
  );
}
import "./Hero.css";

import Typewriter from "typewriter-effect";
import GitHubIcon from "@mui/icons-material/GitHub";

const Hero = () => {
  return (
    <div id="home" className="hero">
      <div className="content">
        <h1 className="center-text">
          Vibe Check
          <a
            href="https://github.com/machine-moon/vibecheck"
            className="social-icons social-icon github"
            target="_blank"
            rel="noopener noreferrer"
          >
            <GitHubIcon />
          </a>
        </h1>

        <div className="center-text">
          <Typewriter
            options={{
              autoStart: true,
              loop: true,
              delay: 60,
              strings: [
                "Welcome to Vibe Check 🚀",
                "Manage tweets and solve tweet-based problems with hints.",
              ],
              pauseFor: 1200,
            }}
          />
        </div>
      </div>
    </div>
  );
};

export default Hero;

import "./Footer.css";
import GitHubIcon from "@mui/icons-material/GitHub";

const Footer = () => {
  const currentYear = new Date().getFullYear();

  return (
    <div className="footer">
      <hr />
      <div className="footer-bottom">
        <p className="footer-bottom-left" style={{ color: "black" }}>
          © {currentYear} Tarek Ibrahim
        </p>
        <div className="social-icons">
          <a
            href="https://github.com/machine-moon/"
            className="social-icon github"
            target="_blank"
            rel="noopener noreferrer"
          >
            <GitHubIcon />
          </a>
        </div>
      </div>
    </div>
  );
};

export default Footer;

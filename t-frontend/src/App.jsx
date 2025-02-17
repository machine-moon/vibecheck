import Hero from "./Components/Hero/Hero.jsx";
import Contact from "./Components/Contact/Contact.jsx";
import Footer from "./Components/Footer/Footer.jsx";
import VibeCheck from "./Components/vibecheck/vibecheck.jsx";

const App = () => {
  return (
    <div className="container">
      <Hero />
      <VibeCheck />
      <Contact />
      <Footer />
    </div>
  );
};

export default App;

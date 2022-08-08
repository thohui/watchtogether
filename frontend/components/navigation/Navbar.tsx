export const Navbar = () => {
  return (
    <div className="navbar mb-2 shadow-lg bg-neutral  rounded-box">
      <button
        onClick={() => window.location.replace("/")}
        className="btn btn-ghost"
      >
        <span className="text-xl font-bold">📺 Watchtogether</span>
      </button>
    </div>
  );
};

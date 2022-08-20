export const Navbar = () => {
  return (
    <div className="navbar bg-base-300 mb-2 rounded-lg">
      <button
        onClick={() => window.location.replace("/")}
        className="btn btn-ghost btn-sm"
      >
        <span className="lg:text-xl sm:text-xs">Watchtogether</span>
      </button>
    </div>
  );
};

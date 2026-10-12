import "./Header.css";
import api from "../api";
import { useAuth } from "../auth/AuthContext";


function Header() {
    const { user, setUser } = useAuth();
    const logout = async () => { try { await api.post("/auth/logout"); } finally { setUser(null); } };
    const name = user?.name || user?.Name || "Manager";
    return(
        <header className="header">
            <div className="brand"><span className="brand-mark">G</span><span>Garago</span></div>
            <div className="header-user"><div className="user-avatar">{name.charAt(0).toUpperCase()}</div><div className="user-details"><strong>{name}</strong><span>Manager</span></div><button className="logout-button" onClick={logout}>Log out</button></div>
        </header>
    );
}

export default Header;

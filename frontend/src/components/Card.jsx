
function Card({ title, value, caption, tone = "neutral", icon }) {
    return (
        <div className={`card metric-card tone-${tone}`}>
            <div className="metric-heading">{icon && <span className="metric-icon" aria-hidden="true">{icon}</span>}<h3>{title}</h3></div>
            <p>{value}</p>
            {caption && <span className="card-caption">{caption}</span>}
        </div>
    );
}

export default Card;

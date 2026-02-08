import { ShoppingCart, CreditCard, Package } from "lucide-react";
import "./index.css";

const NavBar = ({ onCart, onCheckout, onOrders }) => {
  return (
    <nav className="navbar">
      <div className="navbar-left">
        {["A", "B", "C", "D", "E"].map((l) => (
          <div key={l} className="logo-circle">{l}</div>
        ))}
      </div>

      <div className="navbar-right">
        <button className="nav-btn" onClick={onCart}>
          <ShoppingCart size={20} />
          <span>Cart</span>
        </button>

        <button className="nav-btn" onClick={onCheckout}>
          <CreditCard size={20} />
          <span>Checkout</span>
        </button>

        <button className="nav-btn" onClick={onOrders}>
          <Package size={20} />
          <span>Orders</span>
        </button>
      </div>
    </nav>
  );
};

export default NavBar;

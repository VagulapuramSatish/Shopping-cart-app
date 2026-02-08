import React, { useEffect, useState } from "react";
import api from "../../api";
import "./index.css";

const ItemsList = ({ token, action, clearAction }) => {
  const [items, setItems] = useState([]);
  const [cartItems, setCartItems] = useState([]);
  const [orders, setOrders] = useState([]);

  useEffect(() => {
    fetchItems();
  }, []);

  useEffect(() => {
    if (!action) return;

    if (action === "cart") showCart();
    else if (action === "checkout") checkout();
    else if (action === "orders") showOrders();

    clearAction();
  }, [action]);

  
  const fetchItems = async () => {
    try {
      const res = await api.get("/items");
      setItems(res.data);
    } catch (err) {
      console.error(err);
      alert("Failed to load items");
    }
  };


  const fetchCart = async () => {
    try {
      const res = await api.get("/carts", {
        headers: { Authorization: "Bearer " + token },
      });
      setCartItems(res.data || []);
    } catch (err) {
      console.error(err);
      setCartItems([]);
    }
  };

  
  const fetchOrders = async () => {
    try {
      const res = await api.get("/orders", {
        headers: { Authorization: "Bearer " + token },
      });
      setOrders(res.data || []);
    } catch (err) {
      console.error(err);
      setOrders([]);
    }
  };

  
  const addToCart = async (item) => {
    try {
      await api.post(
        "/carts",
        { item_ids: [item.ID] },
        { headers: { Authorization: "Bearer " + token } }
      );
      await fetchCart();
      alert(`Added "${item.Name}" to cart!`);
    } catch (err) {
      console.error(err);
      alert("Failed to add item to cart");
    }
  };

  
  const showCart = async () => {
    await fetchCart();
    if (!cartItems.length) return alert("Cart is empty");

    const msg = cartItems
      .map((i) => `CartID: ${i.CartID}, ItemID: ${i.ItemID}`)
      .join("\n");
    alert(msg);
  };

 
  const checkout = async () => {
    await fetchCart();
    if (!cartItems.length) return alert("Cart is empty");

    try {
      await api.post(
        "/orders",
        { cart_id: cartItems[0].CartID },
        { headers: { Authorization: "Bearer " + token } }
      );
      setCartItems([]);
      alert("Order successful!");
      fetchOrders();
    } catch (err) {
      console.error(err);
      alert("Checkout failed");
    }
  };

 
  const showOrders = async () => {
    await fetchOrders();
    if (!orders.length) return alert("No orders yet");

    const msg = orders.map((o) => `OrderID: ${o.ID}`).join("\n");
    alert(msg);
  };

  return (
    <div className="items-page">
      <h1 className="items-title">Items</h1>

      <div className="items-grid">
        {items.map((item) => (
          <div key={item.ID} className="item-card" onClick={() => addToCart(item)}>
            <div className="item-image">{/* Placeholder for future images */}</div>
            <h3 className="item-name">{item.Name}</h3>
            <p className="item-status">Status: {item.Status}</p>
          </div>
        ))}
      </div>
    </div>
  );
};

export default ItemsList;

"use client";

import React, { useEffect, useState } from "react";
import { Form, Button, Alert } from "react-bootstrap";

type コミュニティ = {
  ID: string;
  Name: string;
};

export default function InterferencePage() {
  // コミュニティ一覧を保持
  const [communities, setCommunities] = useState<コミュニティ[]>([]);
  // 選択されたコミュニティID
  const [commA, setCommA] = useState("");
  const [commB, setCommB] = useState("");
  // 実行結果メッセージ
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  // コミュニティ一覧を取得
  async function fetchCommunities() {
    try {
      const res = await fetch("http://localhost:8080/communities");
      if (!res.ok) {
        throw new Error("Failed to fetch communities");
      }
      const data = await res.json();
      setCommunities(data);
    } catch (err: any) {
      console.error(err);
      setError("Error fetching communities.");
    }
  }

  // 干渉シミュレーション実行
  async function handleInterferenceSim() {
    setMessage("");
    setError("");
    if (!commA || !commB) {
      setError("Please select two different communities.");
      return;
    }
    if (commA === commB) {
      setError("コミュニティ A and B must be different.");
      return;
    }
    try {
      const url = `http://localhost:8080/simulate/interference?commA=${commA}&commB=${commB}`;
      const res = await fetch(url, { method: "POST" });
      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || "Failed to simulate interference.");
      }
      const result = await res.json();
      setMessage(`Interference simulation done. ${JSON.stringify(result)}`);
    } catch (err: any) {
      console.error(err);
      setError(err.message || "Unknown error during interference simulation.");
    }
  }

  // マウント時にコミュニティ一覧を取得
  useEffect(() => {
    fetchCommunities();
  }, []);

  return (
    <div>
      <h2>干渉シュミレーション</h2>
      {error && <Alert variant="danger">{error}</Alert>}
      {message && <Alert variant="success">{message}</Alert>}

      <Form>
        <Form.Group className="mb-3" controlId="commA">
          <Form.Label>コミュニティ A</Form.Label>
          <Form.Select value={commA} onChange={(e) => setCommA(e.target.value)}>
            <option value="">-- 選択 --</option>
            {communities.map((c) => (
              <option key={c.ID} value={c.ID}>
                {c.Name}
              </option>
            ))}
          </Form.Select>
        </Form.Group>

        <Form.Group className="mb-3" controlId="commB">
          <Form.Label>コミュニティ B</Form.Label>
          <Form.Select value={commB} onChange={(e) => setCommB(e.target.value)}>
            <option value="">-- 選択 --</option>
            {communities.map((c) => (
              <option key={c.ID} value={c.ID}>
                {c.Name}
              </option>
            ))}
          </Form.Select>
        </Form.Group>

        <Button variant="primary" onClick={handleInterferenceSim}>
          実行
        </Button>
      </Form>
    </div>
  );
}

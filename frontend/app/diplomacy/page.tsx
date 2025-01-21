"use client";

import React, { useEffect, useState } from "react";
import { Form, Button } from "react-bootstrap";

type Community = {
  ID: string;
  Name: string;
};

export default function DiplomacyPage() {
  const [communities, setCommunities] = useState<Community[]>([]);
  const [commA, setCommA] = useState("");
  const [commB, setCommB] = useState("");

  // コミュニティ一覧を取得
  async function fetchCommunities() {
    try {
      const res = await fetch("http://localhost:8080/communities");
      if (!res.ok) {
        throw new Error("Failed to fetch communities");
      }
      const data = await res.json();
      setCommunities(data);
    } catch (err) {
      console.error(err);
    }
  }

  // 外交シミュレーション実行
  async function handleDiplomacySim() {
    if (!commA || !commB || commA === commB) {
      alert("Please select two different communities.");
      return;
    }
    try {
      const res = await fetch(
        `http://localhost:8080/simulate/diplomacy?commA=${commA}&commB=${commB}`,
        { method: "POST" }
      );
      if (!res.ok) {
        throw new Error("Diplomacy simulation failed");
      }
      alert("Diplomacy simulation done.");
    } catch (err) {
      console.error(err);
      alert("Error occurred during diplomacy simulation.");
    }
  }

  // マウント時にコミュニティ一覧を取得
  useEffect(() => {
    fetchCommunities();
  }, []);

  return (
    <div>
      <h2>Diplomacy Simulation</h2>

      <Form>
        {/* Community A */}
        <Form.Group className="mb-3" controlId="communityA">
          <Form.Label>Community A</Form.Label>
          <Form.Select value={commA} onChange={(e) => setCommA(e.target.value)}>
            <option value="">-- select --</option>
            {communities.map((c) => (
              <option key={c.ID} value={c.ID}>
                {c.Name}
              </option>
            ))}
          </Form.Select>
        </Form.Group>

        {/* Community B */}
        <Form.Group className="mb-3" controlId="communityB">
          <Form.Label>Community B</Form.Label>
          <Form.Select value={commB} onChange={(e) => setCommB(e.target.value)}>
            <option value="">-- select --</option>
            {communities.map((c) => (
              <option key={c.ID} value={c.ID}>
                {c.Name}
              </option>
            ))}
          </Form.Select>
        </Form.Group>

        {/* 実行ボタン */}
        <Button variant="primary" onClick={handleDiplomacySim}>
          Simulate Diplomacy
        </Button>
      </Form>
    </div>
  );
}

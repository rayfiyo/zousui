"use client";

import React, { useEffect, useState } from "react";
import { Form, Button } from "react-bootstrap";

type コミュニティ = {
  ID: string;
  Name: string;
};

export default function DiplomacyPage() {
  const [communities, setCommunities] = useState<コミュニティ[]>([]);
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
      <h2>外交シミュレーション</h2>

      <Form>
        {/* コミュニティ A */}
        <Form.Group className="mb-3" controlId="communityA">
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

        {/* コミュニティ B */}
        <Form.Group className="mb-3" controlId="communityB">
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

        <Button variant="primary" onClick={handleDiplomacySim}>
          実行
        </Button>
      </Form>
    </div>
  );
}

'use client';

import React, { useEffect, useState } from 'react';

type Community = {
  ID: string;
  Name: string;
};

export default function DiplomacyPage() {
  const [communities, setCommunities] = useState<Community[]>([]);
  const [commA, setCommA] = useState('');
  const [commB, setCommB] = useState('');

  // コミュニティ一覧取得
  async function fetchCommunities() {
    try {
      const res = await fetch('http://localhost:8080/communities');
      if (!res.ok) {
        throw new Error('Failed to fetch communities');
      }
      const data = await res.json();
      setCommunities(data);
    } catch (err) {
      console.error(err);
    }
  }

  async function handleDiplomacySim() {
    if (!commA || !commB || commA === commB) {
      alert('Please select two different communities.');
      return;
    }
    try {
      const res = await fetch(
        `http://localhost:8080/simulate/diplomacy?commA=${commA}&commB=${commB}`,
        {
          method: 'POST',
        }
      );
      if (!res.ok) {
        throw new Error('Diplomacy simulation failed');
      }
      alert('Diplomacy simulation done.');
    } catch (err) {
      console.error(err);
      alert('Error occurred during diplomacy simulation.');
    }
  }

  useEffect(() => {
    fetchCommunities();
  }, []);

  return (
    <div>
      <h2>Diplomacy Simulation</h2>
      <div className="mb-3">
        <label className="form-label">Community A</label>
        <select
          className="form-select"
          value={commA}
          onChange={(e) => setCommA(e.target.value)}
        >
          <option value="">-- select --</option>
          {communities.map((c) => (
            <option key={c.ID} value={c.ID}>
              {c.Name}
            </option>
          ))}
        </select>
      </div>
      <div className="mb-3">
        <label className="form-label">Community B</label>
        <select
          className="form-select"
          value={commB}
          onChange={(e) => setCommB(e.target.value)}
        >
          <option value="">-- select --</option>
          {communities.map((c) => (
            <option key={c.ID} value={c.ID}>
              {c.Name}
            </option>
          ))}
        </select>
      </div>
      <button className="btn btn-primary" onClick={handleDiplomacySim}>
        Simulate Diplomacy
      </button>
    </div>
  );
}

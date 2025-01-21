// app/community/[id]/page.tsx
"use client";

import React, { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";

type Community = {
  ID: string;
  Name: string;
  Description: string;
  Population: number;
  Culture: string;
};

export default function CommunityDetailPage() {
  const params = useParams();
  const router = useRouter();
  const communityID = params.id as string;

  const [community, setCommunity] = useState<Community | null>(null);

  async function fetchCommunity() {
    try {
      const res = await fetch(
        `http://localhost:8080/communities/${communityID}`
      );
      if (!res.ok) {
        throw new Error(`Error: ${res.status}`);
      }
      const data = await res.json();
      setCommunity(data);
    } catch (err) {
      console.error("Error fetching community:", err);
    }
  }

  // 更新をしたい場合(省略可能)
  async function handleEdit() {
    // 例: NameやDescriptionを変更するフォームを表示し、PATCH or PUTするフロー
    alert("Edit functionality is not implemented yet.");
  }

  // シミュレーション呼び出し
  async function handleSimulate() {
    if (!community) return;
    try {
      const res = await fetch(
        `http://localhost:8080/simulate/${community.ID}`,
        {
          method: "POST",
        }
      );
      if (!res.ok) {
        throw new Error("Simulation failed.");
      }
      alert("Simulated successfully.");
      fetchCommunity(); // 再取得
    } catch (err) {
      console.error(err);
    }
  }

  useEffect(() => {
    if (communityID) {
      fetchCommunity();
    }
  }, [communityID]);

  if (!community) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <h2>Community Detail: {community.Name}</h2>
      <p>ID: {community.ID}</p>
      <p>Description: {community.Description}</p>
      <p>Population: {community.Population}</p>
      <p>Culture: {community.Culture}</p>

      <button className="btn btn-primary me-2" onClick={handleSimulate}>
        Simulate
      </button>
      <button className="btn btn-secondary me-2" onClick={handleEdit}>
        Edit (TODO)
      </button>
      <button className="btn btn-light" onClick={() => router.push("/")}>
        Back
      </button>
    </div>
  );
}

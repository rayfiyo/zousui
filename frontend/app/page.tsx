'use client';

import React, { useEffect, useState } from 'react';
import Link from 'next/link';

type Community = {
  ID: string;
  Name: string;
  Population: number;
  Culture: string;
};

export default function HomePage() {
  const [communities, setCommunities] = useState<Community[]>([]);

  // 1. コミュニティ一覧を取得
  async function fetchCommunities() {
    try {
      const res = await fetch('http://localhost:8080/communities');
      if (!res.ok) {
        throw new Error(`Error: ${res.status} ${res.statusText}`);
      }
      const data = await res.json();
      setCommunities(data);
    } catch (err) {
      console.error('Error fetching communities:', err);
    }
  }

  // 2. コミュニティ削除
  async function handleDelete(id: string) {
    const confirmed = window.confirm(
      'Are you sure you want to delete this community?'
    );
    if (!confirmed) return;
    try {
      const res = await fetch(`http://localhost:8080/communities/${id}`, {
        method: 'DELETE',
      });
      if (!res.ok) {
        alert(`Delete failed. Status: ${res.status}`);
        return;
      }
      alert('Community deleted.');
      fetchCommunities();
    } catch (err) {
      console.error('Error deleting community:', err);
    }
  }

  // 3. シミュレーションAPI呼び出し
  async function handleSimulate(communityID: string) {
    try {
      const res = await fetch(`http://localhost:8080/simulate/${communityID}`, {
        method: 'POST',
      });
      if (!res.ok) {
        alert(`Simulation failed. Status: ${res.status}`);
        return;
      }
      alert('Simulation executed successfully.');

      // 再取得して表示更新
      fetchCommunities();
    } catch (err) {
      console.error('Error simulating:', err);
    }
  }

  // 初回マウント時にコミュニティ一覧を取得
  useEffect(() => {
    fetchCommunities();
  }, []);

  return (
    <main>
      {/* <main className="container mt-4">*/}
      <h1 className="mb-4">Zousui Communities</h1>
      {communities.length === 0 ? (
        <p>No communities found.</p>
      ) : (
        <div className="row">
          {communities.map((comm) => (
            <div key={comm.ID} className="col-md-4 mb-3">
              <div className="card">
                <div className="card-body">
                  <h5 className="card-title">{comm.Name}</h5>
                  <p className="card-text">
                    Population: {comm.Population} <br />
                    Culture: {comm.Culture}
                  </p>
                  {/* シミュレーション */}
                  <button
                    className="btn btn-primary me-2"
                    onClick={() => handleSimulate(comm.ID)}
                  >
                    Simulate
                  </button>
                  {/* 詳細ページへのリンク */}
                  <Link
                    href={`/community/${comm.ID}`}
                    className="btn btn-secondary me-2"
                  >
                    Details
                  </Link>
                  {/* 削除 */}
                  <button
                    className="btn btn-danger"
                    onClick={() => handleDelete(comm.ID)}
                  >
                    Delete
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </main>
  );
}

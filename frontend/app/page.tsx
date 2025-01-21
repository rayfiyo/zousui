"use client";

import React, { useEffect, useState } from 'react';

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

    // 2. シミュレーションAPI呼び出し
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
        <main className="container mt-4">
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
                                    <button
                                        className="btn btn-primary"
                                        onClick={() => handleSimulate(comm.ID)}
                                    >
                                        Simulate
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

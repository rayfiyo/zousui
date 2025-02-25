"use client";

import React, { useEffect, useState } from "react";
import Link from "next/link";
import { Button, Row, Col, Card } from "react-bootstrap";

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
      const res = await fetch("http://localhost:8080/communities");
      if (!res.ok) {
        throw new Error(`Error: ${res.status} ${res.statusText}`);
      }
      const data = await res.json();
      setCommunities(data);
    } catch (err) {
      console.error("Error fetching communities:", err);
    }
  }

  // 2. コミュニティ削除
  async function handleDelete(id: string) {
    const confirmed = window.confirm(
      "Are you sure you want to delete this community?"
    );
    if (!confirmed) return;
    try {
      const res = await fetch(`http://localhost:8080/communities/${id}`, {
        method: "DELETE",
      });
      if (!res.ok) {
        alert(`Delete failed. Status: ${res.status}`);
        return;
      }
      alert("Community deleted.");
      fetchCommunities();
    } catch (err) {
      console.error("Error deleting community:", err);
    }
  }

  // 3. シミュレーションAPI呼び出し
  async function handleSimulate(communityID: string) {
    try {
      const res = await fetch(`http://localhost:8080/simulate/${communityID}`, {
        method: "POST",
      });
      if (!res.ok) {
        alert(`Simulation failed. Status: ${res.status}`);
        return;
      }
      alert("Simulation executed successfully.");

      // 再取得して表示更新
      fetchCommunities();
    } catch (err) {
      console.error("Error simulating:", err);
    }
  }

  // 初回マウント時にコミュニティ一覧を取得
  useEffect(() => {
    fetchCommunities();
  }, []);

  return (
    <main>
      <h1 className="my-4 text-center">zousui Communities</h1>

      {communities.length === 0 ? (
        <p>No communities found.</p>
      ) : (
        <Row>
          {communities.map((comm) => (
            <Col key={comm.ID} xs={12} className="mb-3">
              <Card>
                <Card.Body>
                  <Card.Title>{comm.Name}</Card.Title>
                  <Card.Text>
                    Population: {comm.Population} <br />
                    Culture: {comm.Culture}
                  </Card.Text>
                  {/* ボタングループ（縦並び）を中央寄せ */}
                  <div className="d-flex justify-content-center">
                    {/* シミュレーション */}
                    <div className="d-flex flex-wrap gap-2">
                      <Button
                        variant="primary"
                        onClick={() => handleSimulate(comm.ID)}
                      >
                        シュミレート
                      </Button>

                      {/* 詳細ページへのリンク（Next.js Link を as で指定） */}
                      <Link href={`/community/${comm.ID}`}>
                        <Button variant="secondary">詳細</Button>
                      </Link>

                      {/* 削除 */}
                      <Button
                        variant="danger"
                        onClick={() => handleDelete(comm.ID)}
                      >
                        削除
                      </Button>
                    </div>
                  </div>
                </Card.Body>
              </Card>
            </Col>
          ))}
        </Row>
      )}
    </main>
  );
}

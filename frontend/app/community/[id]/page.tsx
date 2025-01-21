"use client";

import React, { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { Button, Card, Spinner, Container, Row, Col } from "react-bootstrap";

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
  const [loading, setLoading] = useState(false);

  async function fetchCommunity() {
    try {
      setLoading(true);
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
    } finally {
      setLoading(false);
    }
  }

  // 更新機能（サンプル）
  async function handleEdit() {
    alert("Edit functionality is not implemented yet.");
  }

  // シミュレーション
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

  // 読み込み中にスピナーを表示
  if (loading) {
    return (
      <div className="text-center mt-5">
        <Spinner animation="border" role="status">
          <span className="visually-hidden">Loading community detail...</span>
        </Spinner>
      </div>
    );
  }

  if (!community) {
    return <div>No community data found.</div>;
  }

  return (
    <Container>
      <Row className="justify-content-center">
        <Col md={8}>
          <h2 className="my-4">Community Detail</h2>
          <Card>
            <Card.Body>
              <Card.Title>{community.Name}</Card.Title>
              <Card.Text>ID: {community.ID}</Card.Text>
              <Card.Text>Description: {community.Description}</Card.Text>
              <Card.Text>Population: {community.Population}</Card.Text>
              <Card.Text>Culture: {community.Culture}</Card.Text>

              <div className="d-flex gap-2 mt-4">
                <Button variant="primary" onClick={handleSimulate}>
                  Simulate
                </Button>
                <Button variant="secondary" onClick={handleEdit}>
                  Edit (TODO)
                </Button>
                <Button variant="light" onClick={() => router.push("/")}>
                  Back
                </Button>
              </div>
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </Container>
  );
}

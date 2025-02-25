"use client";

import React, { useEffect, useState, useCallback } from "react";
import { useParams, useRouter } from "next/navigation";
import Image from "next/image";
import {
  Button,
  Card,
  Spinner,
  Container,
  Row,
  Col,
  Alert,
} from "react-bootstrap";

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

  // 画像生成関連の state
  const [imageSrc, setImageSrc] = useState<string>("");
  const [imageError, setImageError] = useState<string>("");

  // ====== Fetch Community ====== //
  const fetchCommunity = useCallback(async () => {
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
  }, [communityID]);

  // ====== Delete ====== //
  async function handleDelete() {
    if (!community) return;
    if (!confirm("Are you sure you want to delete this community?")) {
      return;
    }

    try {
      const res = await fetch(
        `http://localhost:8080/communities/${community.ID}`,
        {
          method: "DELETE",
        }
      );
      if (!res.ok) {
        alert(`Delete failed. Status: ${res.status}`);
        return;
      }
      alert("Community deleted.");
      router.push("/"); // 一覧に戻る
    } catch (err) {
      console.error("Error deleting community:", err);
    }
  }

  // ====== Simulate ====== //
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

  // ====== Interference Simulation ====== //
  async function handleInterference() {
    if (!community) return;
    try {
      const res = await fetch(
        `http://localhost:8080/simulate/interference/${community.ID}`,
        {
          method: "POST",
        }
      );
      if (!res.ok) {
        throw new Error("Interference simulation failed.");
      }
      alert("Interference simulation succeeded.");
      fetchCommunity(); // 最新データを再取得
    } catch (err) {
      console.error(err);
      alert("Error occurred during interference simulation.");
    }
  }

  // ====== Edit (stub) ======
  async function handleEdit() {
    alert("Edit functionality is not implemented yet.");
  }

  // ====== Generate Image ======
  async function handleGenerateImage() {
    if (!community) return;

    // 一度リセット
    setImageSrc("");
    setImageError("");

    try {
      setLoading(true);
      // 例: "POST /communities/:communityID/generateImage"
      const url = `http://localhost:8080/communities/${community.ID}/generateImage`;

      // もし style などを追加したい場合は body JSON を入れる
      // const body = JSON.stringify({ style: "fantasy" });
      // const headers = { "Content-Type": "application/json" };

      const res = await fetch(url, {
        method: "POST",
        // headers,
        // body,
      });

      if (!res.ok) {
        throw new Error(`Image generation request failed: ${res.statusText}`);
      }

      // 画像は blob で受け取る
      const blob = await res.blob();
      const objectUrl = URL.createObjectURL(blob);
      setImageSrc(objectUrl);
    } catch (err: unknown) {
      if (err instanceof Error) {
        setImageError(err.message);
      } else {
        setImageError("Unknown error");
      }
      console.error("Error generating image:", err);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    if (communityID) {
      fetchCommunity();
    }
  }, [fetchCommunity]);

  // 読み込み中にスピナーを表示
  if (loading && !community) {
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
          <h2 className="my-4 text-center">Community Detail</h2>
          <Card>
            <Card.Body>
              <Card.Title className="text-center">{community.Name}</Card.Title>
              <Card.Text>ID: {community.ID}</Card.Text>
              <Card.Text>Description: {community.Description}</Card.Text>
              <Card.Text>Population: {community.Population}</Card.Text>
              <Card.Title className="text-center"> Culture </Card.Title>
              <Card.Text> {community.Culture}</Card.Text>

              <div className="d-flex flex-wrap gap-2 mt-4 justify-content-center">
                <Button variant="primary" onClick={handleSimulate}>
                  シュミレート
                </Button>
                <Button variant="warning" onClick={handleInterference}>
                  他国との干渉
                </Button>

                <Button variant="success" onClick={handleGenerateImage}>
                  画像生成
                </Button>
              </div>

              <div className="d-flex flex-wrap gap-2 mt-4 justify-content-center">
                <Button variant="danger" onClick={handleDelete}>
                  削除
                </Button>
                <Button variant="light" onClick={() => router.push("/")}>
                  一覧に戻る
                </Button>
                <Button variant="secondary" onClick={handleEdit}>
                  編集
                </Button>
              </div>

              {loading && community && (
                <div className="mt-3 text-center">
                  <Spinner animation="border" role="status">
                    <span className="visually-hidden">Generating image...</span>
                  </Spinner>
                </div>
              )}

              {imageError && (
                <Alert variant="danger" className="mt-3">
                  Error: {imageError}
                </Alert>
              )}

              {imageSrc && (
                <div className="mt-3">
                  <h5>Generated Image:</h5>
                  <Image
                    src={imageSrc}
                    alt="Generated"
                    style={{ maxWidth: "100%" }}
                  />
                </div>
              )}
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </Container>
  );
}

"use client";

import React, { useEffect, useState } from "react";
import { Alert, Button, Card, Col, Form, Row, Spinner } from "react-bootstrap";

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
  // ユーザーが追加入力するフィールド
  const [userInput, setUserInput] = useState("");
  // 実行結果メッセージ
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const [isSimulating, setIsSimulating] = useState(false);
  const [simulationResult, setSimulationResult] = useState<any>(null);

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
    setError("");
    setMessage("");
    setSimulationResult(null);
    if (!commA || !commB) {
      setError("Please select two different communities.");
      return;
    }
    if (commA === commB) {
      setError("コミュニティ A and B must be different.");
      return;
    }
    setIsSimulating(true);
    try {
      const payload = { commA, commB, userInput };
      const res = await fetch("http://localhost:8080/simulate/interference", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || "Failed to simulate interference.");
      }
      const result = await res.json();
      setMessage("Interference simulation completed successfully.");
      setSimulationResult(result);
    } catch (err: any) {
      setError(err.message || "Unknown error during interference simulation.");
    } finally {
      setIsSimulating(false);
    }
  }

  useEffect(() => {
    fetchCommunities();
  }, []);

  return (
    <div>
      <h2 className="mb-4 text-center">干渉シミュレーション</h2>
      <Card className="mb-4">
        <Card.Body>
          <Card.Title>干渉シミュレーションについて</Card.Title>
          <Card.Text>
            このシミュレーションでは、2つの異なるコミュニティを選択し、
            外部からの影響がどのように影響し合うかを見ることができるよ。
            あと、集計プロセスに影響を与えるために、
            ユーザー自身の入力を提供することもできるよ。
          </Card.Text>
        </Card.Body>
      </Card>

      {error && <Alert variant="danger">{error}</Alert>}
      {message && <Alert variant="success">{message}</Alert>}

      <Form>
        <Row className="mb-3">
          <Col md={4}>
            <Form.Group controlId="communityA">
              <Form.Label>コミュニティ A</Form.Label>
              <Form.Select
                value={commA}
                onChange={(e) => setCommA(e.target.value)}
              >
                <option value="">-- コミュニティ A を選択 --</option>
                {communities.map((c) => (
                  <option key={c.ID} value={c.ID}>
                    {c.Name}
                  </option>
                ))}
              </Form.Select>
            </Form.Group>
          </Col>
          <Col md={4}>
            <Form.Group controlId="communityB">
              <Form.Label>コミュニティ B</Form.Label>
              <Form.Select
                value={commB}
                onChange={(e) => setCommB(e.target.value)}
              >
                <option value="">-- コミュニティ B を選択 --</option>
                {communities.map((c) => (
                  <option key={c.ID} value={c.ID}>
                    {c.Name}
                  </option>
                ))}
              </Form.Select>
            </Form.Group>
          </Col>
          <Col md={4}>
            <Form.Group controlId="userInput">
              <Form.Label>追加入力</Form.Label>
              <Form.Control
                as="textarea"
                rows={3}
                placeholder="その他のアイデアやインスピレーションがあれば、ここに"
                value={userInput}
                onChange={(e) => setUserInput(e.target.value)}
              />
            </Form.Group>
          </Col>
        </Row>
        <div className="d-flex justify-content-center mb-4">
          <Button
            variant="primary"
            onClick={handleInterferenceSim}
            disabled={isSimulating}
          >
            {isSimulating ? (
              <>
                <Spinner animation="border" size="sm" /> シュミレート中...
              </>
            ) : (
              "実行"
            )}
          </Button>
        </div>
      </Form>

      {simulationResult && (
        <Card className="mb-4">
          <Card.Body>
            <Card.Title>Simulation Result</Card.Title>
            <pre>{JSON.stringify(simulationResult, null, 2)}</pre>
          </Card.Body>
        </Card>
      )}
    </div>
  );
}

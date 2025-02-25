"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { Form, Button } from "react-bootstrap";

export default function NewCommunityPage() {
  const router = useRouter();

  const [formData, setFormData] = useState({
    id: "",
    name: "",
    description: "",
    population: 0,
    culture: "",
  });

  function handleChange(
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    // バリデーション省略
    try {
      const res = await fetch("http://localhost:8080/communities", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          id: formData.id,
          name: formData.name,
          description: formData.description,
          population: Number(formData.population),
          culture: formData.culture,
        }),
      });
      if (!res.ok) {
        throw new Error(`Failed to create community. Status: ${res.status}`);
      }
      alert("Community created successfully!");
      router.push("/");
    } catch (err) {
      console.error(err);
      alert("Error creating community");
    }
  }

  return (
    <div>
      <h2>Create New Community</h2>

      <Form onSubmit={handleSubmit} className="mt-3">
        {/* ID */}
        <Form.Group className="mb-3" controlId="community-id">
          <Form.Label>ID</Form.Label>
          <Form.Control
            type="text"
            name="id"
            value={formData.id}
            onChange={handleChange}
            placeholder="comm-123..."
            required
          />
        </Form.Group>

        {/* Name */}
        <Form.Group className="mb-3" controlId="community-name">
          <Form.Label>Name</Form.Label>
          <Form.Control
            type="text"
            name="name"
            value={formData.name}
            onChange={handleChange}
            required
          />
        </Form.Group>

        {/* Description */}
        <Form.Group className="mb-3" controlId="community-description">
          <Form.Label>Description</Form.Label>
          <Form.Control
            as="textarea"
            rows={3}
            name="description"
            value={formData.description}
            onChange={handleChange}
          />
        </Form.Group>

        {/* Population */}
        <Form.Group className="mb-3" controlId="community-population">
          <Form.Label>Population</Form.Label>
          <Form.Control
            type="number"
            name="population"
            value={formData.population}
            onChange={handleChange}
          />
        </Form.Group>

        {/* Culture */}
        <Form.Group className="mb-3" controlId="community-culture">
          <Form.Label>Culture</Form.Label>
          <Form.Control
            as="textarea"
            rows={2}
            name="culture"
            value={formData.culture}
            onChange={handleChange}
          />
        </Form.Group>

        {/* Submit button */}
        <Button variant="primary" type="submit">
          Create
        </Button>
      </Form>
    </div>
  );
}

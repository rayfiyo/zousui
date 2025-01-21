'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';

export default function NewCommunityPage() {
  const router = useRouter();

  const [formData, setFormData] = useState({
    id: '',
    name: '',
    description: '',
    population: 0,
    culture: '',
  });

  function handleChange(
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    // バリデーション省略。idを自動生成にしたいなら不要
    try {
      const res = await fetch('http://localhost:8080/communities', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
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
      alert('Community created successfully!');
      router.push('/'); // 一覧に戻る
    } catch (err) {
      console.error(err);
      alert('Error creating community');
    }
  }

  return (
    <div>
      <h2>Create New Community</h2>
      <form onSubmit={handleSubmit} className="mt-3">
        <div className="mb-3">
          <label className="form-label">ID</label>
          <input
            className="form-control"
            type="text"
            name="id"
            value={formData.id}
            onChange={handleChange}
            placeholder="comm-123..."
            required
          />
        </div>
        <div className="mb-3">
          <label className="form-label">Name</label>
          <input
            className="form-control"
            type="text"
            name="name"
            value={formData.name}
            onChange={handleChange}
            required
          />
        </div>
        <div className="mb-3">
          <label className="form-label">Description</label>
          <textarea
            className="form-control"
            name="description"
            value={formData.description}
            onChange={handleChange}
          />
        </div>
        <div className="mb-3">
          <label className="form-label">Population</label>
          <input
            className="form-control"
            type="number"
            name="population"
            value={formData.population}
            onChange={handleChange}
          />
        </div>
        <div className="mb-3">
          <label className="form-label">Culture</label>
          <textarea
            className="form-control"
            name="culture"
            value={formData.culture}
            onChange={handleChange}
          />
        </div>
        <button type="submit" className="btn btn-primary">
          Create
        </button>
      </form>
    </div>
  );
}

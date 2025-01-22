"use client";
import React, { useState } from 'react';

export default function CommunityDetailPage({ params }: { params: { id: string } }) {
  const [imgSrc, setImgSrc] = useState<string>("");

  async function handleGenerateImage() {
    const url = `http://localhost:8080/communities/${params.id}/generateImage`;

    // このままだと直接IMGではなく fetch でバイナリ取得→URL化する例
    try {
      const res = await fetch(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        // bodyに style など追加してもOK
        // body: JSON.stringify({ style: "horror" })
      });
      if (!res.ok) {
        alert("Error generating image");
        return;
      }
      const blob = await res.blob();
      const objectUrl = URL.createObjectURL(blob);
      setImgSrc(objectUrl);
    } catch (err) {
      console.error(err);
    }
  }

  return (
    <div>
      <h2>Community {params.id}</h2>
      {/* ...コミュニティ情報表示... */}
      <button onClick={handleGenerateImage}>Generate Image</button>

      {imgSrc && (
        <div>
          <h3>Generated Image:</h3>
          <img src={imgSrc} alt="Generated" />
        </div>
      )}
    </div>
  );
}


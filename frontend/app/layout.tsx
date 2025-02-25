import { ReactNode } from "react";
import NavBar from "./NavBar";

export const metadata = {
  title: "Zousui Communities",
  description: "文明進化シミュレーター / Civilization Evolution Simulator",
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="ja">
      <body>
        {/* クライアントコンポーネントの呼び出し */}
        <NavBar />

        <div className="container mt-4">{children}</div>
      </body>
    </html>
  );
}

import 'bootstrap/dist/css/bootstrap.min.css';
import { ReactNode } from 'react';
import Link from 'next/link';

export const metadata = {
  title: 'Zousui Communities',
  description: '文明進化シミュレーター / Civilization Evolution Simulator',
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="ja">
      <body>
        {/* ナビバー */}
        <nav className="navbar navbar-expand-lg navbar-dark bg-dark">
          <div className="container">
            <Link className="navbar-brand" href="/">
              Zousui
            </Link>
            <button
              className="navbar-toggler"
              type="button"
              data-bs-toggle="collapse"
              data-bs-target="#navbarNav"
              aria-controls="navbarNav"
              aria-expanded="false"
              aria-label="Toggle navigation"
            >
              <span className="navbar-toggler-icon"></span>
            </button>
            <div className="collapse navbar-collapse" id="navbarNav">
              <ul className="navbar-nav ms-auto">
                <li className="nav-item">
                  <Link className="nav-link" href="/">
                    Home
                  </Link>
                </li>
                <li className="nav-item">
                  <Link className="nav-link" href="/community/new">
                    Create Community
                  </Link>
                </li>
                <li className="nav-item">
                  <Link className="nav-link" href="/diplomacy">
                    Diplomacy
                  </Link>
                </li>
              </ul>
            </div>
          </div>
        </nav>

        {/* ページコンテンツ */}
        <div className="container mt-4">{children}</div>
      </body>
    </html>
  );
}

"use client";

import "bootstrap/dist/css/bootstrap.min.css";
import Link from "next/link";
import { Navbar, Nav, Container } from "react-bootstrap";

export default function NavBar() {
  return (
    <Navbar bg="dark" variant="dark" expand="lg">
      <Container>
        {/* ロゴ */}
        <Navbar.Brand as={Link} href="/">
          zousui
        </Navbar.Brand>

        {/* ハンバーガーメニュー */}
        <Navbar.Toggle aria-controls="navbarNav" />

        {/* ナビゲーションメニュー */}
        <Navbar.Collapse id="navbarNav">
          <Nav className="ms-auto">
            <Nav.Link as={Link} href="/">
              ホーム
            </Nav.Link>
            <Nav.Link as={Link} href="/community/new">
              新規作成
            </Nav.Link>
            <Nav.Link as={Link} href="/diplomacy">
              外交
            </Nav.Link>
            <Nav.Link as={Link} href="/interference">
              干渉
            </Nav.Link>
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
}

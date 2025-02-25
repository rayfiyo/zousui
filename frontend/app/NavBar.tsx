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
          Zousui
        </Navbar.Brand>

        {/* ハンバーガーメニュー */}
        <Navbar.Toggle aria-controls="navbarNav" />

        {/* ナビゲーションメニュー */}
        <Navbar.Collapse id="navbarNav">
          <Nav className="ms-auto">
            <Nav.Link as={Link} href="/">
              Home
            </Nav.Link>
            <Nav.Link as={Link} href="/community/new">
              Create Community
            </Nav.Link>
            <Nav.Link as={Link} href="/diplomacy">
              Diplomacy
            </Nav.Link>
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
}

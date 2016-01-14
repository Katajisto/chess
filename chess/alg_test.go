package chess

import (
	"strings"
	"testing"
)

type algDecodeTest struct {
	Pos     *Position
	AlgText string
	PostPos *Position
}

type algEncodeTest struct {
	Pos     *Position
	Move    *Move
	AlgText string
}

var (
	validAlgDecodeTests = []algDecodeTest{
		{
			// opening for white
			Pos:     unsafeFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"),
			AlgText: "e4",
			PostPos: unsafeFEN("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"),
		},
		{
			// http://en.lichess.org/W91M4jms#14
			Pos:     unsafeFEN("rn1qkb1r/pp3ppp/2p1pn2/3p4/2PP4/2NQPN2/PP3PPP/R1B1K2R b KQkq - 0 7"),
			AlgText: "Nbd7",
			PostPos: unsafeFEN("r2qkb1r/pp1n1ppp/2p1pn2/3p4/2PP4/2NQPN2/PP3PPP/R1B1K2R w KQkq - 1 8"),
		},
		{
			// http://en.lichess.org/W91M4jms#17
			Pos:     unsafeFEN("r2qk2r/pp1n1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P4PPP/R1B1K2R w KQkq - 1 9"),
			AlgText: "O-O",
			PostPos: unsafeFEN("r2qk2r/pp1n1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P4PPP/R1B2RK1 b kq - 0 9"),
		},
		{
			// http://en.lichess.org/W91M4jms#23
			Pos:     unsafeFEN("3r1rk1/pp1nqppp/2pbpn2/3p4/2PP4/1PNQPN2/PB3PPP/3RR1K1 b - - 5 12"),
			AlgText: "dxc4",
			PostPos: unsafeFEN("3r1rk1/pp1nqppp/2pbpn2/8/2pP4/1PNQPN2/PB3PPP/3RR1K1 w - - 0 13"),
		},
		{
			// http://en.lichess.org/W91M4jms#34
			Pos:     unsafeFEN("3r2k1/pp2qppp/2p2n2/4b3/2P5/2N1P3/PB2QPPP/3R2K1 w - - 1 18"),
			AlgText: "Rxd8+",
			PostPos: unsafeFEN("3R2k1/pp2qppp/2p2n2/4b3/2P5/2N1P3/PB2QPPP/6K1 b - - 0 18"),
		},
		{
			// http://en.lichess.org/W91M4jms#2
			Pos:     unsafeFEN("rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR b KQkq d3 0 1"),
			AlgText: "Nf6",
			PostPos: unsafeFEN("rnbqkb1r/pppppppp/5n2/8/3P4/8/PPP1PPPP/RNBQKBNR w KQkq - 1 2"),
		},
		{
			// http://en.lichess.org/4HXJOtpN#26
			Pos:     unsafeFEN("r4rk1/ppbn1p2/2p1pn1p/q2p2pb/7N/PP1PP1PP/1BPN1PB1/R3QRK1 w - g6 0 14"),
			AlgText: "Nhf3",
			PostPos: unsafeFEN("r4rk1/ppbn1p2/2p1pn1p/q2p2pb/8/PP1PPNPP/1BPN1PB1/R3QRK1 b - - 1 14"),
		},
		{
			// http://en.lichess.org/4HXJOtpN#87
			Pos:     unsafeFEN("4r3/8/2p2PPk/1p1r4/pP2p1R1/P1B5/2P2K2/8 b - - 0 44"),
			AlgText: "Rd1??",
			PostPos: unsafeFEN("4r3/8/2p2PPk/1p6/pP2p1R1/P1B5/2P2K2/3r4 w - - 1 45"),
		},
		{
			// http://en.lichess.org/YXPuk8kg#79
			Pos:     unsafeFEN("8/3k2Kp/p7/1p1r4/6P1/8/Pn1p4/7R b - - 0 40"),
			AlgText: "d1=Q",
			PostPos: unsafeFEN("8/3k2Kp/p7/1p1r4/6P1/8/Pn6/3q3R w - - 0 41"),
		},
		{
			// http://en.lichess.org/dimuEVR0#18
			Pos:     unsafeFEN("rnbk1b1r/p3pppp/5n2/2p1p3/5B2/2N2P2/PPP3PP/R3KBNR w KQ - 0 10"),
			AlgText: "O-O-O+",
			PostPos: unsafeFEN("rnbk1b1r/p3pppp/5n2/2p1p3/5B2/2N2P2/PPP3PP/2KR1BNR b - - 0 10"),
		},
		{
			// only 1 rook can move because of pin http://en.lichess.org/JCRBhXH7#62
			Pos:     unsafeFEN("4R3/1r1k2pp/p1p5/1pP5/8/8/1PP3PP/2K1Rr2 w - - 5 32"),
			AlgText: "Re7+",
			PostPos: unsafeFEN("8/1r1kR1pp/p1p5/1pP5/8/8/1PP3PP/2K1Rr2 b - - 6 32"),
		},
	}

	invalidAlgDecodeTests = []algDecodeTest{
		{
			// opening for white
			Pos:     unsafeFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"),
			AlgText: "e5",
		},
		{
			// http://en.lichess.org/W91M4jms#14
			Pos:     unsafeFEN("rn1qkb1r/pp3ppp/2p1pn2/3p4/2PP4/2NQPN2/PP3PPP/R1B1K2R b KQkq - 0 7"),
			AlgText: "Nd7",
		},
		{
			// http://en.lichess.org/W91M4jms#17
			Pos:     unsafeFEN("r2qk2r/pp1n1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P4PPP/R1B1K2R w KQkq - 1 9"),
			AlgText: "O-O-O-O",
			PostPos: unsafeFEN("r2qk2r/pp1n1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P4PPP/R1B2RK1 b kq - 0 9"),
		},
		{
			// http://en.lichess.org/W91M4jms#23
			Pos:     unsafeFEN("3r1rk1/pp1nqppp/2pbpn2/3p4/2PP4/1PNQPN2/PB3PPP/3RR1K1 b - - 5 12"),
			AlgText: "dx4",
		},
		{
			// should not assume pawn for unknown peice type "n"
			Pos:     unsafeFEN("rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2"),
			AlgText: "nf3",
		},
	}

	algEncodeTests = []algEncodeTest{
		{
			Pos:     unsafeFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"),
			Move:    &Move{s1: E2, s2: E4},
			AlgText: "e4",
		},
		{
			Pos:     unsafeFEN("rn1qkb1r/pp3ppp/2p1pn2/3p4/2PP4/2NbPN2/PP3PPP/R1BQK2R w KQkq - 0 7"),
			Move:    &Move{s1: D1, s2: D3},
			AlgText: "Qxd3",
		},
		{
			Pos:     unsafeFEN("rn1qkb1r/pp3ppp/2p1pn2/3p4/2PP4/2NQPN2/PP3PPP/R1B1K2R b KQkq - 0 7"),
			Move:    &Move{s1: B8, s2: D7},
			AlgText: "Nbd7",
		},
		{
			Pos:     unsafeFEN("r2qk2r/pp1n1ppp/2pbpn2/3p4/2PP4/1PNQPN2/P4PPP/R1B1K2R w KQkq - 1 9"),
			Move:    &Move{s1: E1, s2: G1},
			AlgText: "O-O",
		},
		{
			Pos:     unsafeFEN("r4rk1/pp1nqppp/2pbpn2/3p4/2PP4/1PNQPN2/PB3PPP/R4RK1 w - - 2 11"),
			Move:    &Move{s1: A1, s2: D1},
			AlgText: "Rad1",
		},
		{
			Pos:     unsafeFEN("3q2k1/pp3ppp/2p2n2/4b3/2P5/2N1P3/PB3PPP/3Q2K1 b - - 1 19"),
			Move:    &Move{s1: D8, s2: D1},
			AlgText: "Qxd1+",
		},
		{
			Pos:     unsafeFEN("rnbk1b1r/p3pppp/5n2/2p1p3/5B2/2N2P2/PPP3PP/R3KBNR w KQ - 0 10"),
			Move:    &Move{s1: E1, s2: C1},
			AlgText: "O-O-O+",
		},
		{
			Pos:     unsafeFEN("8/8/8/4p3/3P4/8/8/8 w - - 0 10"),
			Move:    &Move{s1: D4, s2: E5},
			AlgText: "dxe5",
		},
		{
			Pos:     unsafeFEN("8/8/8/4p3/3P1P2/8/8/8 w - - 0 10"),
			Move:    &Move{s1: F4, s2: E5},
			AlgText: "fxe5",
		},
	}
)

func TestValidAlgDecoding(t *testing.T) {
	for _, test := range validAlgDecodeTests {
		move, err := decodeMove(test.Pos, test.AlgText)
		if err != nil {
			movesStrList := []string{}
			for _, m := range test.Pos.ValidMoves() {
				movesStrList = append(movesStrList, encodeMove(test.Pos, m))
			}
			t.Fatalf("starting from board \n%s\n expected move to be valid error - %s %s\n", test.Pos.board.Draw(), err, strings.Join(movesStrList, ","))
		}
		postPos := test.Pos.Update(move)
		if postPos.String() != test.PostPos.String() {
			t.Fatalf("starting from board \n%s\n after move %s\n expected board to be %s\n%s\n but was %s\n%s\n",
				test.Pos.board.Draw(), move.String(), test.PostPos.String(),
				test.PostPos.board.Draw(), postPos.String(), postPos.board.Draw())
		}
	}
}

//
// func TestInvalidAlgDecoding(t *testing.T) {
// 	for _, test := range invalidAlgDecodeTests {
// 		if _, err := decodeMove(test.PreState, test.AlgText); err == nil {
// 			t.Fatalf("starting from board\n%s\n expected move notation %s to be invalid", test.PreState.board.Draw(), test.AlgText)
// 		}
// 	}
// }
//
// func TestValidAlgEncoding(t *testing.T) {
// 	for _, test := range algEncodeTests {
// 		if !test.Move.isValid() {
// 			t.Fatalf("starting from board\n%s\n invalid move %s", test.Move.state.board.Draw(), test.Move.String())
// 		}
// 		actual := encodeMove(test.Move)
// 		if actual != test.AlgText {
// 			t.Fatalf("starting from board\n%s\n to board \n%s\nexpected move notation %s but got %s", test.Move.state.board.Draw(), test.Move.PostMoveState().board.Draw(), test.AlgText, actual)
// 		}
// 	}
// }

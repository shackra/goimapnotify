package main

// This file is part of goimapnotify
// Copyright (C) 2017  Jorge Javier Araya Navarro

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

import (
	"sync"
)

type guardian struct {
	wg      *sync.WaitGroup
	senders int
	mx      *sync.Mutex
}

func (g *guardian) Add(delta int) {
	g.wg.Add(delta)
	g.mx.Lock()
	g.senders = g.senders + delta
	g.mx.Unlock()
}

func (g *guardian) Done() {
	g.wg.Done()
	g.mx.Lock()
	if g.senders > 0 {
		g.senders = g.senders - 1
	}
	g.mx.Unlock()
}

func (g *guardian) Wait() {
	g.wg.Wait()
}

func (g *guardian) Close(channel chan<- IDLEEvent) {
	g.mx.Lock()
	if g.senders == 1 {
		close(channel)
	}
	g.mx.Unlock()
}

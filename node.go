package main

import (
	"image"
	"time"

	"github.com/clktmr/n64/drivers/controller"
)

type Spatial interface {
	Position() image.Point
	setGlobal(image.Point)
}

type Node struct {
	relativePos, globalPos image.Point

	children []Updater
	parent   Updater
}

func NewNode(children ...Updater) *Node {
	node := &Node{children: children}
	for _, child := range children {
		child.setParent(node)
	}
	return node
}

func (p *Node) Update(delta time.Duration, input [4]controller.Controller) {
	for _, child := range p.children {
		if child, ok := child.(Spatial); ok {
			child.setGlobal(p.globalPos.Add(child.Position()))
		}
		child.Update(delta, input)
	}
}

func (p *Node) Position() image.Point {
	return p.relativePos
}

func (p *Node) setGlobal(pos image.Point) {
	p.globalPos = pos
}

func (p *Node) Children() []Updater {
	return []Updater(p.children)
}

func (p *Node) AddChild(child Updater) {
	p.children = append(p.children, child)
	if child, ok := child.(Spatial); ok {
		child.setGlobal(p.globalPos.Add(child.Position()))
	}
	child.setParent(p)
}

func (p *Node) setParent(parent Updater) {
	p.parent = parent
}

func (p *Node) RemoveChild(child Updater) {
	for i, c := range p.children {
		if c == child {
			c.setParent(nil)
			p.children = append(p.children[:i], p.children[i+1:]...)
		}
	}
}

package gru

import (
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strconv"

	G "gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

type layer struct {
	wf *tensor.Dense
	uf *tensor.Dense
	bf *tensor.Dense

	wh *tensor.Dense
	uh *tensor.Dense
	bh *tensor.Dense

	ones *tensor.Dense
}

// Model is a GRU model
type Model struct {
	layers []*layer
	we     *tensor.Dense
	be     *tensor.Dense
	wo     *tensor.Dense
	bo     *tensor.Dense

	inputs                               int
	inputSize, embeddingSize, outputSize int
	layerSizes                           []int
}

// NewModel creates a new GRU model
func NewModel(rnd *rand.Rand, inputs, inputSize, embeddingSize, outputSize int, layerSizes []int) *Model {
	gaussian32 := func(s ...int) []float32 {
		size := tensor.Shape(s).TotalSize()
		weights, stdev := make([]float32, size), math.Sqrt(2/float64(s[len(s)-1]))
		for i := range weights {
			weights[i] = float32(rnd.NormFloat64() * stdev)
		}
		return weights
	}

	model := &Model{
		inputs:        inputs,
		inputSize:     inputSize,
		embeddingSize: embeddingSize,
		outputSize:    outputSize,
		layerSizes:    layerSizes,
	}
	model.we = tensor.New(tensor.WithShape(embeddingSize, inputSize),
		tensor.WithBacking(gaussian32(embeddingSize, inputSize)))
	model.be = tensor.New(tensor.Of(tensor.Float32), tensor.WithShape(embeddingSize))

	previous := inputs * embeddingSize
	for _, size := range layerSizes {
		layer := &layer{}
		model.layers = append(model.layers, layer)

		layer.wf = tensor.New(tensor.WithShape(size, previous),
			tensor.WithBacking(gaussian32(size, previous)))
		layer.uf = tensor.New(tensor.WithShape(size, size),
			tensor.WithBacking(gaussian32(size, size)))
		layer.bf = tensor.New(tensor.Of(tensor.Float32), tensor.WithShape(size))

		layer.wh = tensor.New(tensor.WithShape(size, previous),
			tensor.WithBacking(gaussian32(size, previous)))
		layer.uh = tensor.New(tensor.WithShape(size, size),
			tensor.WithBacking(gaussian32(size, size)))
		layer.bh = tensor.New(tensor.Of(tensor.Float32), tensor.WithShape(size))

		layer.ones = tensor.Ones(tensor.Float32, size)

		previous = size
	}

	model.wo = tensor.New(tensor.WithShape(outputSize, previous),
		tensor.WithBacking(gaussian32(outputSize, previous)))
	model.bo = tensor.New(tensor.Of(tensor.Float32), tensor.WithShape(outputSize))

	return model
}

// WriteFile writes the weights to a file
func (m *Model) WriteFile(file string) error {
	out, err := os.Create(file)
	if err != nil {
		return err
	}
	defer out.Close()
	return m.Write(out)
}

// Write writes the weights to a Writer
func (m *Model) Write(out io.Writer) error {
	encoder := gob.NewEncoder(out)
	write := func(t *tensor.Dense) error {
		return encoder.Encode(t.Data())
	}

	for _, layer := range m.layers {
		err := write(layer.wf)
		if err != nil {
			return err
		}
		err = write(layer.uf)
		if err != nil {
			return err
		}
		err = write(layer.bf)
		if err != nil {
			return err
		}

		err = write(layer.wh)
		if err != nil {
			return err
		}
		err = write(layer.uh)
		if err != nil {
			return err
		}
		err = write(layer.bh)
		if err != nil {
			return err
		}
	}
	err := write(m.we)
	if err != nil {
		return err
	}
	err = write(m.be)
	if err != nil {
		return err
	}
	err = write(m.wo)
	if err != nil {
		return err
	}
	err = write(m.bo)
	if err != nil {
		return err
	}

	return nil
}

// Read reads the weights from a Reader
func (m *Model) Read(in io.Reader) error {
	decoder := gob.NewDecoder(in)
	read := func(t *tensor.Dense) error {
		data := t.Data().([]float32)
		return decoder.Decode(&data)
	}

	for _, layer := range m.layers {
		err := read(layer.wf)
		if err != nil {
			return err
		}
		err = read(layer.uf)
		if err != nil {
			return err
		}
		err = read(layer.bf)
		if err != nil {
			return err
		}

		err = read(layer.wh)
		if err != nil {
			return err
		}
		err = read(layer.uh)
		if err != nil {
			return err
		}
		err = read(layer.bh)
		if err != nil {
			return err
		}
	}
	err := read(m.we)
	if err != nil {
		return err
	}
	err = read(m.be)
	if err != nil {
		return err
	}
	err = read(m.wo)
	if err != nil {
		return err
	}
	err = read(m.bo)
	if err != nil {
		return err
	}

	return nil
}

// ReadFile reads the weights from a file
func (m *Model) ReadFile(file string) error {
	in, err := os.Open(file)
	if err != nil {
		return nil
	}
	defer in.Close()
	return m.Read(in)
}

func (m *Model) compare(b *Model) error {
	compare := func(a, b *tensor.Dense, name string) error {
		x := a.Data().([]float32)
		y := b.Data().([]float32)
		for k, v := range x {
			if v != y[k] {
				return fmt.Errorf("%v %v %v %v %v", k, name, "they don't match", v, y[k])
			}
		}
		return nil
	}
	for i, layer := range m.layers {
		err := compare(layer.wf, b.layers[i].wf, "wf"+string(i))
		if err != nil {
			return err
		}
		err = compare(layer.uf, b.layers[i].uf, "uf"+string(i))
		if err != nil {
			return err
		}
		err = compare(layer.bf, b.layers[i].bf, "bf"+string(i))
		if err != nil {
			return err
		}
		err = compare(layer.wh, b.layers[i].wh, "wh"+string(i))
		if err != nil {
			return err
		}
		err = compare(layer.uh, b.layers[i].uh, "uh"+string(i))
		if err != nil {
			return err
		}
		err = compare(layer.bh, b.layers[i].bh, "bh"+string(i))
		if err != nil {
			return err
		}
	}
	err := compare(m.we, b.we, "we")
	if err != nil {
		return err
	}
	err = compare(m.be, b.be, "be")
	if err != nil {
		return err
	}
	err = compare(m.wo, b.wo, "wo")
	if err != nil {
		return err
	}
	err = compare(m.bo, b.bo, "bo")
	if err != nil {
		return err
	}

	return nil
}

type gru struct {
	wf *G.Node
	uf *G.Node
	bf *G.Node

	wh *G.Node
	uh *G.Node
	bh *G.Node

	ones *G.Node
}

func (l *layer) NewGRULayer(g *G.ExprGraph, name string) *gru {
	wf := G.NodeFromAny(g, l.wf, G.WithName("wf_"+name))
	uf := G.NodeFromAny(g, l.uf, G.WithName("uf_"+name))
	bf := G.NodeFromAny(g, l.bf, G.WithName("bf_"+name))

	wh := G.NodeFromAny(g, l.wh, G.WithName("wh_"+name))
	uh := G.NodeFromAny(g, l.uh, G.WithName("uh_"+name))
	bh := G.NodeFromAny(g, l.bh, G.WithName("bh_"+name))

	ones := G.NodeFromAny(g, l.ones, G.WithName("ones_"+name))
	return &gru{
		wf:   wf,
		uf:   uf,
		bf:   bf,
		wh:   wh,
		uh:   uh,
		bh:   bh,
		ones: ones,
	}
}

func (g *gru) fwd(input, previous *G.Node) *G.Node {
	x := G.Must(G.Mul(g.wf, input))
	y := G.Must(G.Mul(g.uf, previous))
	f := G.Must(G.Sigmoid(G.Must(G.Add(G.Must(G.Add(x, y)), g.bf))))

	x = G.Must(G.Mul(g.wh, input))
	y = G.Must(G.Mul(g.uh, G.Must(G.HadamardProd(f, previous))))
	z := G.Must(G.Tanh(G.Must(G.Add(G.Must(G.Add(x, y)), g.bh))))

	a := G.Must(G.HadamardProd(G.Must(G.Sub(g.ones, f)), z))
	b := G.Must(G.HadamardProd(f, previous))

	return G.Must(G.Add(a, b))
}

type gruOut struct {
	hiddens       G.Nodes
	probabilities *G.Node
}

// RNN is a LSTM that takes characters as input
type RNN struct {
	*Model
	layers []*gru

	g       *G.ExprGraph
	we      *G.Node
	be      *G.Node
	wo      *G.Node
	bo      *G.Node
	hiddens G.Nodes

	steps    int
	inputs   [][]*tensor.Dense
	outputs  []*tensor.Dense
	previous []*gruOut
	cost     *G.Node
	machine  G.VM
}

// NewRNN create a new GRU for characters as inputs
func NewRNN(model *Model) *RNN {
	g := G.NewGraph()
	var layers []*gru
	var hiddens G.Nodes
	for i, v := range model.layerSizes {
		name := strconv.Itoa(i)
		layer := model.layers[i].NewGRULayer(g, name)
		layers = append(layers, layer)

		hiddenTensor := tensor.New(tensor.Of(tensor.Float32), tensor.WithShape(v))
		hidden := G.NewVector(g, G.Float32, G.WithName("prevHidden_"+name),
			G.WithShape(v), G.WithValue(hiddenTensor))
		hiddens = append(hiddens, hidden)
	}
	we := G.NodeFromAny(g, model.we, G.WithName("we"))
	be := G.NodeFromAny(g, model.be, G.WithName("be"))
	wo := G.NodeFromAny(g, model.wo, G.WithName("wo"))
	bo := G.NodeFromAny(g, model.bo, G.WithName("bo"))
	return &RNN{
		Model:   model,
		layers:  layers,
		g:       g,
		we:      we,
		be:      be,
		wo:      wo,
		bo:      bo,
		hiddens: hiddens,
	}
}

func (r *RNN) learnables() (value G.Nodes) {
	for _, l := range r.layers {
		nodes := G.Nodes{
			l.wf,
			l.uf,
			l.bf,
			l.wh,
			l.uh,
			l.bh,
		}
		value = append(value, nodes...)
	}

	value = append(value, r.we)
	value = append(value, r.be)
	value = append(value, r.wo)
	value = append(value, r.bo)

	return
}

func (r *RNN) fwd(previous *gruOut) (inputs []*tensor.Dense, retVal *gruOut, err error) {
	previousHiddens := r.hiddens
	if previous != nil {
		previousHiddens = previous.hiddens
	}

	var hiddens G.Nodes
	for i, v := range r.layers {
		var inputVector *G.Node
		if i == 0 {
			inputs = make([]*tensor.Dense, r.Model.inputs)
			for j := range inputs {
				inputs[j] = tensor.New(tensor.Of(tensor.Float32), tensor.WithShape(r.inputSize))
				input := G.NewVector(r.g, tensor.Float32, G.WithShape(r.inputSize), G.WithValue(inputs[j]))
				if inputVector == nil {
					inputVector = G.Must(G.Add(G.Must(G.Mul(r.we, input)), r.be))
				} else {
					inputVector = G.Must(G.Concat(0, inputVector, G.Must(G.Add(G.Must(G.Mul(r.we, input)), r.be))))
				}
			}
		} else {
			inputVector = hiddens[i-1]
		}

		hidden := v.fwd(inputVector, previousHiddens[i])
		hiddens = append(hiddens, hidden)
	}
	lastHidden := hiddens[len(hiddens)-1]
	var output *G.Node
	if output, err = G.Mul(r.wo, lastHidden); err == nil {
		if output, err = G.Add(output, r.bo); err != nil {
			G.WithName("LAST HIDDEN")(lastHidden)
			ioutil.WriteFile("err.dot", []byte(lastHidden.RestrictedToDot(3, 10)), 0644)
			panic(fmt.Sprintf("ERROR: %v", err))
		}
	} else {
		panic(err)
	}

	var probs *G.Node
	probs = G.Must(G.SoftMax(output))

	retVal = &gruOut{
		hiddens:       hiddens,
		probabilities: probs,
	}

	return
}

func (r *RNN) feedback(tap int) {
	prev := r.previous[tap]
	for i := range r.hiddens {
		input := r.hiddens[i].Value().(*tensor.Dense)
		output := prev.hiddens[i].Value().(*tensor.Dense)
		err := output.CopyTo(input)
		if err != nil {
			panic(err)
		}
	}
}

func (r *RNN) reset() {
	for i := range r.hiddens {
		r.hiddens[i].Value().(*tensor.Dense).Zero()
	}
}

// ModeLearn puts the CharRNN into a learning mode
func (r *RNN) ModeLearn(steps int) (err error) {
	inputs := make([][]*tensor.Dense, r.Model.inputs)
	outputs := make([]*tensor.Dense, steps)
	previous := make([]*gruOut, steps)
	var cost *G.Node

	for i := range inputs {
		inputs[i] = make([]*tensor.Dense, steps)
	}

	for i := 0; i < steps; i++ {
		var loss *G.Node

		var prev *gruOut
		if i > 0 {
			prev = previous[i-1]
		}
		var in []*tensor.Dense
		in, previous[i], err = r.fwd(prev)
		if err != nil {
			return
		}
		for k, v := range in {
			inputs[k][i] = v
		}

		logprob := G.Must(G.Neg(G.Must(G.Log(previous[i].probabilities))))
		outputs[i] = tensor.New(tensor.Of(tensor.Float32), tensor.WithShape(r.outputSize))
		output := G.NewVector(r.g, tensor.Float32, G.WithShape(r.outputSize), G.WithValue(outputs[i]))
		loss = G.Must(G.Mul(logprob, output))

		if cost == nil {
			cost = loss
		} else {
			cost = G.Must(G.Add(cost, loss))
		}
		G.WithName("Cost")(cost)
	}

	r.steps = steps
	r.inputs = inputs
	r.outputs = outputs
	r.previous = previous
	r.cost = cost

	_, err = G.Grad(cost, r.learnables()...)
	if err != nil {
		return
	}

	r.machine = G.NewTapeMachine(r.g, G.BindDualValues(r.learnables()...))
	return
}

// ModeInference puts the CharRNN into inference mode
func (r *RNN) ModeInference() (err error) {
	inputs := make([][]*tensor.Dense, r.Model.inputs)
	previous := make([]*gruOut, 1)

	for i := range inputs {
		inputs[i] = make([]*tensor.Dense, 1)
	}

	var in []*tensor.Dense
	in, previous[0], err = r.fwd(nil)
	if err != nil {
		return
	}
	for k, v := range in {
		inputs[k][0] = v
	}

	r.inputs = inputs
	r.previous = previous
	r.machine = G.NewTapeMachine(r.g)
	return
}

func (r *RNN) getProbabilities(input []int) G.Value {
	end := len(input) - 1
	r.reset()
	for i := range input {
		r.inputs[0][0].Zero()
		r.inputs[0][0].SetF32(input[i], 1.0)
		if len(r.inputs) > 1 {
			r.inputs[1][0].Zero()
			r.inputs[1][0].SetF32(input[end-i], 1.0)
		}
		err := r.machine.RunAll()
		if err != nil {
			panic(err)
		}
		r.feedback(0)
		r.machine.Reset()
	}

	return r.previous[0].probabilities.Value()
}

// AttackProbability return the probability the input is an attack
func (r *RNN) AttackProbability(input []int) (float32, error) {
	value := r.getProbabilities(input)
	if t, ok := value.(tensor.Tensor); ok {
		isAttack, err := t.At(0)
		if err != nil {
			return 0, err
		}
		probability, ok := isAttack.(float32)
		if !ok {
			return 0, fmt.Errorf("value is not float32")
		}
		return 100 * probability, nil
	}

	return 0, fmt.Errorf("not a tensor")
}

// IsAttack determines if an input is an attack
func (r *RNN) IsAttack(input []int) bool {
	value := r.getProbabilities(input)
	if t, ok := value.(tensor.Tensor); ok {
		max, err := tensor.Argmax(t, -1)
		if err != nil {
			panic(err)
		}
		if !max.IsScalar() {
			panic("expected scalar index")
		}
		if x := max.ScalarValue().(int); x == 0 {
			return true
		}
	} else {
		panic("not a tensor")
	}

	return false
}

// Learn learns strings
func (r *RNN) Learn(data []int, attack bool, solver G.Solver) (retCost, retPerp []float64, err error) {
	end := len(data) - 1

	r.reset()
	for i := range data[:len(data)-r.steps+1] {
		for j := 0; j < r.steps; j++ {
			index := i + j
			source, rsource := data[index], data[end-index]

			r.inputs[0][j].Zero()
			r.inputs[0][j].SetF32(source, 1.0)
			if len(r.inputs) > 1 {
				r.inputs[1][j].Zero()
				r.inputs[1][j].SetF32(rsource, 1.0)
			}
			if r.outputs != nil {
				r.outputs[j].Zero()
				if attack {
					r.outputs[j].SetF32(0, 1.0)
				} else {
					r.outputs[j].SetF32(1, 1.0)
				}
			}
		}

		err = r.machine.RunAll()
		if err != nil {
			return
		}

		err = solver.Step(r.learnables())
		if err != nil {
			return
		}

		if cv, ok := r.cost.Value().(G.Scalar); ok {
			retCost = append(retCost, float64(cv.Data().(float32)))
		}
		r.feedback(0)
		r.machine.Reset()
	}

	return
}

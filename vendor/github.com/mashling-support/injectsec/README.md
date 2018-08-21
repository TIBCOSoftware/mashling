[![godoc](https://godoc.org/github.com/pointlander/injectsec?status.svg)](https://godoc.org/github.com/pointlander/injectsec)

# injectsec_train options
```
Usage of injectsec_train:
  -chunks
    	generate chunks
  -data string
    	use data for training
  -epochs int
    	the number of epochs for training (default 1)
  -help
    	print help
  -print
    	print training data
```

# usage of injectsec_train to train a model
```
injectsec_train -data training_data_example.csv --epochs 10
```

Will train using the builtin data set and training_data_example.csv for 10 epochs. The output weights will be placed in a directory named 'output'.

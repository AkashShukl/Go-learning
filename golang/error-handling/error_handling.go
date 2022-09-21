package erratum

import (
	"errors"
	"fmt"
)

func aux(resource Resource, input string) (err error, ok1 bool) {
	defer func() {
		if r := recover(); r != nil {
			nr, ok := r.(FrobError)
			if ok {
				// fmt.Println("reaching", nr, nr.defrobTag)
				resource.Defrob(nr.defrobTag)
				ok1 = true
				err = nr
			}
			_, true := r.(error)
			if true {
				// resource.Close()
				err = r.(error)
				ok1 = true
			}
		}
	}()
	resource.Frob(input)
	return err, ok1
}

func Use(o ResourceOpener, input string) error {
	resource, err := o()
	if err != nil {
		var transientErr TransientError
		// fmt.Printf("%T", err)
		if errors.As(err, &transientErr) {
			for err != nil && errors.As(err, &transientErr) {
				resource, err = o()
			}
		} else if err != nil {
			return err
		}
	} else {
		defer resource.Close()
	}

	err11, ok11 := aux(resource, input)
	// fmt.Println("=>>>", err11, ok11)
	if ok11 {
		err = err11
	}
	fmt.Println("=>>>", err)
	return err
}

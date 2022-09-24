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
				resource.Defrob(nr.defrobTag)
				ok1 = true
				err = nr
			}
			newErrr, true := r.(error)
			err = newErrr
			ok1 = true
		}
	}()
	resource.Frob(input)
	return err, ok1
}

func Use(o ResourceOpener, input string) error {
	resource, err := o()
	if err != nil {
		var transientErr TransientError
		if errors.As(err, &transientErr) {
			for err != nil && errors.As(err, &transientErr) {
				resource, err = o()
			}
		}
		if err != nil {
			return err
		}
	}
	fmt.Println(resource)
	defer resource.Close()

	err11, ok11 := aux(resource, input)
	if ok11 {
		err = err11
	}
	return err
}

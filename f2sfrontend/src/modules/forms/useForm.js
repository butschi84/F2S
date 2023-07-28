import  _ from '../../@lodash';
import { useCallback, useEffect, useState } from 'react';
import moment from 'moment-timezone';

function isValidDate(d) {
	const a = new Date(d) || null;
	return a instanceof Date && !isNaN(a);
}

function useForm(initialState = {}, onSubmit, validate) {
	const [form, setForm] = useState(initialState);
	const [formErrors, setErrors] = useState({});
	const [formValid, setFormValid] = useState(false);
	const [touchedInputs, setTouchedInputs] = useState([]); 

	useEffect(() => {
		if(!validate || Object.keys(form).length === 0) return;
		const errors = validate(form);
		// only errors for touched for inputs
		const touchedErrors = _.pick(errors, touchedInputs)
		setErrors(touchedErrors);
		setFormValid(Object.keys(errors).length === 0);
	}, [form, touchedInputs, validate]);

	const handleChange = useCallback(event => {
		event.persist();

		setTouchedInputs(_.union(touchedInputs, [event.target.name]));
		setForm(_form => {
			const result = _.setIn(
				{ ..._form },
				event.target.name,
				event.target.type === 'checkbox' ? event.target.checked : event.target.value,
				event.target.name === 'due' && isValidDate(event.target.value) ? moment(event.target.value).tz("Europe/Zurich").format(moment.HTML5_FMT.DATETIME_LOCAL_SECONDS) : event.target.value
			)
			
			return result
		});
	}, [touchedInputs]);

	const resetForm = useCallback(() => {
		if (!_.isEqual(initialState, form)) {
			setForm(initialState);
			setTouchedInputs([]);
			setErrors({});
		}
	}, [form, initialState]);

	const setInForm = useCallback((name, value) => {
		setForm(_form => _.setIn(_form, name, value));
	}, []);

	const handleSubmit = useCallback(
		event => {
			if (event) {
				event.preventDefault();
			}
			// check for form errors
			if (validate) {
				const errors = validate(form);
				setErrors(errors);
				if(Object.keys(errors).length) {
					console.log("form has errors",errors);
					return
				};
			}
			// submit the form
			if (onSubmit) {
				onSubmit();
				setErrors({})
			}
		},
		[onSubmit, form, validate]
	);

	return {
		form,
		handleChange,
		handleSubmit,
		resetForm,
		setForm,
		setInForm,
		formErrors,
		formValid
	};
}

export default useForm;

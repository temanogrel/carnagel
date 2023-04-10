import React from 'react';
import PropTypes from 'prop-types';

export const RenderInput = ({ input, label, type, meta }) => {
  const elementId = 'form_' + Math.random().toString(36).substring(7);

  const { touched, error } = meta;

  return (
    <div className="form-item pb20">
      <label className="pb5" htmlFor={ elementId }>{ label }</label>
      <input { ...input } id={ elementId } placeholder={ label } type={ type }/>
      { touched && error && <p className="text-red pt5">{ error }</p> }
    </div>
  );
};

RenderInput.propTypes = {
  label: PropTypes.string.isRequired,
  type: PropTypes.string.isRequired,
  meta: PropTypes.shape({
    touched: PropTypes.bool.isRequired,
    error: PropTypes.string,
  }),
};

RenderInput.defaultValues = {
  label: ' Missing label',
};

export const RenderCheckbox = ({ input, label, meta }) => {
  const elementId = ' form_' + Math.random().toString(36).substring(7);

  const { touched, error } = meta;

  return (
    <div className=" form-item pb20">
      <div>
        <input { ...input } id={ elementId } type="checkbox"/>
        <label htmlFor={ elementId }>{ label }</label>
        { touched && error && <p className="text-red pt5">{ error }</p> }
      </div>
    </div>
  );
};

RenderCheckbox.propTypes = {
  label: PropTypes.string.isRequired,
  meta: PropTypes.shape({
    touched: PropTypes.bool.isRequired,
    error: PropTypes.string,
  }),
};

RenderCheckbox.defaultValues = {
  label: ' Missing label',
};

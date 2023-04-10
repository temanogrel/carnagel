import React, {Component} from 'react';
import {RecordingImageEntity} from 'api/recording';
import { AnimateUp } from 'utils/animations';

type Props = {
  images: RecordingImageEntity[];
};

export class ImageCollection extends Component<Props> {
  shouldComponentUpdate(nextProps) {
    const idsInNextCollection = nextProps.images.map((image) => image.uuid);

    return this.props.images.length !== nextProps.images.length || this.props.images.some((image) => {
      return idsInNextCollection.indexOf(image.uuid) < 0;
    });
  }

  render() {
    const animateUp = !/PhantomJS/.test(window.navigator.userAgent) ? 'animate-up-opacity-base animate-up-opacity ' : '';

    if (this.props.images.length === 1) {
      const image = this.props.images[0];

      return (
        <div key={image.uuid} className={animateUp + 'grid-12'}>
          <div className='image' style={{paddingBottom: '325%', backgroundSize: 'cover', backgroundImage: image.url}}/>
        </div>
      )
    }

    return (
      <div className='pb30'>
        {this.props.images.map(image => (
          <div key={image.uuid} className={animateUp + 'grid-3 mobile-grid-6'}>
            <div className='image' style={{'backgroundImage': image.url}}/>
          </div>
        ))}
        <div className='clear'/>
      </div>
    );
  }
}

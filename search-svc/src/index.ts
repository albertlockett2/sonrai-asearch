import grpc from 'grpc';
import {
  ISearchServiceServer,
  SearchServiceService
} from './gen/search_grpc_pb';
import {
  Search,
  TransformLegacySearchRequest,
  TransformLegacySearchResponse
} from './gen/search_pb';
import { transformLegacy } from './transform-legacy';

class SearchServer implements ISearchServiceServer {
  transformLegacySearch(
    call: grpc.ServerUnaryCall<TransformLegacySearchRequest>,
    callback: grpc.sendUnaryData<TransformLegacySearchResponse>
  ): void {
    const sourceRaw = call.request.getSource();
    const source = JSON.parse(sourceRaw);
    const search = transformLegacy(source);

    try {
      const response = new TransformLegacySearchResponse();
      response.setSearch(search);
      callback(null, response);
    } catch (e) {
      console.error(e);
      callback(
        {
          name: 'error',
          message: 'error'
        },
        null
      );
    }
  }
}

const server = new grpc.Server();
server.addService<ISearchServiceServer>(
  SearchServiceService,
  new SearchServer()
);
console.log(`Listening on ${process.env.PORT}`);
server.bind(
  `localhost:${process.env.PORT}`,
  grpc.ServerCredentials.createInsecure()
);
server.start();
